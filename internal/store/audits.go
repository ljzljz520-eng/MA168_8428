package store

import (
	"bookstore/recommendation/internal/model"
	bolt "go.etcd.io/bbolt"
	"sort"
)

func (s *Store) SaveAudit(event model.AuditEvent) error {
	data, err := marshal(event)
	if err != nil {
		return err
	}
	return s.withUpdate(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error { return tx.Bucket([]byte("audits")).Put([]byte(event.ID), data) })
	})
}

func (s *Store) ListAudits(recordID string) ([]model.AuditEvent, error) {
	result := make([]model.AuditEvent, 0)
	err := s.withView(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			return tx.Bucket([]byte("audits")).ForEach(func(_, value []byte) error {
				var event model.AuditEvent
				if err := unmarshal(value, &event); err != nil {
					return err
				}
				if recordID == "" || event.RecordID == recordID {
					result = append(result, event)
				}
				return nil
			})
		})
	})
	sort.Slice(result, func(i, j int) bool {
		if result[i].At == result[j].At {
			return result[i].ID < result[j].ID
		}
		return result[i].At < result[j].At
	})
	return result, err
}
