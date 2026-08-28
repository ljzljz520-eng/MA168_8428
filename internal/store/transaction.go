package store

import (
	"bookstore/recommendation/internal/model"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) SaveRecordWithAudit(record model.Record, event model.AuditEvent) error {
	recordData, err := marshal(record)
	if err != nil {
		return err
	}
	eventData, err := marshal(event)
	if err != nil {
		return err
	}
	return s.withUpdate(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error {
			if err := tx.Bucket([]byte("records")).Put([]byte(record.ID), recordData); err != nil {
				return err
			}
			return tx.Bucket([]byte("audits")).Put([]byte(event.ID), eventData)
		})
	})
}

func (s *Store) Snapshot() (map[string]int, error) {
	result := make(map[string]int)
	err := s.withView(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			for _, name := range []string{"records", "audits", "workflows", "attachments"} {
				result[name] = tx.Bucket([]byte(name)).Stats().KeyN
			}
			return nil
		})
	})
	return result, err
}
