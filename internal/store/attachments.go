package store

import (
	"bookstore/recommendation/internal/model"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) SaveAttachment(attachment model.Attachment) error {
	data, err := marshal(attachment)
	if err != nil {
		return err
	}
	return s.withUpdate(func(db *bolt.DB) error {
		return db.Update(func(tx *bolt.Tx) error { return tx.Bucket([]byte("attachments")).Put([]byte(attachment.ID), data) })
	})
}

func (s *Store) GetAttachment(id string) (model.Attachment, error) {
	var attachment model.Attachment
	err := s.withView(func(db *bolt.DB) error {
		return db.View(func(tx *bolt.Tx) error {
			value := tx.Bucket([]byte("attachments")).Get([]byte(id))
			if value == nil {
				return ErrNotFound
			}
			return unmarshal(value, &attachment)
		})
	})
	return attachment, err
}
