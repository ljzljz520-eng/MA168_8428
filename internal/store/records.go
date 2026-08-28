package store

import (
	"sort"

	"bookstore/recommendation/internal/model"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) SaveRecord(record model.Record) error {
	data, err := marshal(record)
	if err != nil {
		return err
	}
	return s.withUpdate(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error { return tx.Bucket([]byte("records")).Put([]byte(record.ID), data) })
	})
}

func (s *Store) GetRecord(id string) (model.Record, error) {
	var record model.Record
	err := s.withView(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			value := tx.Bucket([]byte("records")).Get([]byte(id))
			if value == nil {
				return ErrNotFound
			}
			return unmarshal(value, &record)
		})
	})
	return record, err
}

func (s *Store) ListRecords() ([]model.Record, error) {
	records := make([]model.Record, 0)
	err := s.withView(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte("records")).ForEach(func(_, value []byte) error {
				var record model.Record
				if err := unmarshal(value, &record); err != nil {
					return err
				}
				records = append(records, record)
				return nil
			})
		})
	})
	sort.Slice(records, func(i, j int) bool {
		if records[i].UpdatedAt == records[j].UpdatedAt {
			return records[i].ID < records[j].ID
		}
		return records[i].UpdatedAt < records[j].UpdatedAt
	})
	return records, err
}

func (s *Store) DeleteRecord(id string) error {
	return s.withUpdate(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error { return tx.Bucket([]byte("records")).Delete([]byte(id)) })
	})
}

func (s *Store) CountRecords() (int, error) {
	count := 0
	err := s.withView(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error { count = tx.Bucket([]byte("records")).Stats().KeyN; return nil })
	})
	return count, err
}
