package workflow

import (
	"bookstore/recommendation/internal/model"
	"bookstore/recommendation/internal/recommendation"
	"bookstore/recommendation/internal/store"
	"path/filepath"
	"testing"
)

func TestWorkflowCreateReviewArchive(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "flow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	service := New(st)
	if _, err = service.Register(model.RecordInput{ID: "flow-1", StoreID: "s", Title: "Shelf", Score: 76}, 1, "clerk"); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Submit("flow-1", "clerk", 2); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Review("flow-1", "manager", true, "good", 3); err != nil {
		t.Fatal(err)
	}
	record, err := service.Archive("flow-1", "manager", 4)
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != model.StatusArchived {
		t.Fatalf("archive failed: %s", record.Status)
	}
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "publish.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	service := New(st)
	if _, err = service.Register(model.RecordInput{ID: "pub-1", StoreID: "s", Title: "Catalog", Score: 80}, 1, "clerk"); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Update("pub-1", model.RecordInput{StoreID: "s", Title: "Catalog Revised", Score: 84}, "clerk", 2); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Submit("pub-1", "clerk", 3); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Review("pub-1", "manager", true, "approved", 4); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Publish("pub-1", "manager", 5); err != nil {
		t.Fatal(err)
	}
	collection, err := service.Recommendations(structRequest(), 6)
	if err != nil || len(collection.Items) != 1 {
		t.Fatalf("publish recommendation failed: %v", err)
	}
}

func structRequest() recommendation.Request {
	return recommendation.Request{StoreID: "s", IncludePublished: true}
}
