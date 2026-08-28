package store

import (
	"bookstore/recommendation/internal/model"
	bolt "go.etcd.io/bbolt"
	"sort"
)

func (s *Store) SaveWorkflow(workflow model.Workflow) error {
	data, err := marshal(workflow)
	if err != nil {
		return err
	}
	return s.withUpdate(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error { return tx.Bucket([]byte("workflows")).Put([]byte(workflow.ID), data) })
	})
}

func (s *Store) ListWorkflow(recordID, name string) ([]model.Workflow, error) {
	result := make([]model.Workflow, 0)
	err := s.withView(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte("workflows")).ForEach(func(_, value []byte) error {
				var workflow model.Workflow
				if err := unmarshal(value, &workflow); err != nil {
					return err
				}
				if (recordID == "" || workflow.RecordID == recordID) && (name == "" || workflow.Name == name) {
					result = append(result, workflow)
				}
				return nil
			})
		})
	})
	sort.Slice(result, func(i, j int) bool {
		if result[i].Position == result[j].Position {
			return result[i].ID < result[j].ID
		}
		return result[i].Position < result[j].Position
	})
	return result, err
}

func (s *Store) CompleteWorkflow(id string) error {
	return s.withUpdate(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("workflows"))
			value := bucket.Get([]byte(id))
			if value == nil {
				return ErrNotFound
			}
			var item model.Workflow
			if err := unmarshal(value, &item); err != nil {
				return err
			}
			item.Completed = true
			data, err := marshal(item)
			if err != nil {
				return err
			}
			return bucket.Put([]byte(id), data)
		})
	})
}
