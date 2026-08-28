package store

import (
	"bookstore/recommendation/internal/model"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "records.db")
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := model.Record{ID: "persist-1", StoreID: "store-1", Title: "Long Shelf", Score: 88, Status: model.StatusApproved, Version: 1}
	if err = first.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err = first.SaveAudit(model.NewAuditEvent("audit-1", record.ID, "tester", "approve", 4, model.StatusPending, model.StatusApproved, "ok")); err != nil {
		t.Fatal(err)
	}
	if err = first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	loaded, err := second.GetRecord(record.ID)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Title != record.Title || loaded.Score != record.Score {
		t.Fatalf("reopened record mismatch: %#v", loaded)
	}
	audits, err := second.ListAudits(record.ID)
	if err != nil || len(audits) != 1 {
		t.Fatalf("reopened audit mismatch: %v %d", err, len(audits))
	}
}
