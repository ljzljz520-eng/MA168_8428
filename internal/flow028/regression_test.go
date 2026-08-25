package flow028

import (
	"bookstore/recommendation/internal/model"
	"bookstore/recommendation/internal/recommendation"
	"bookstore/recommendation/internal/store"
	"bookstore/recommendation/internal/workflow"
	"path/filepath"
	"testing"
)

func Test168BusinessRegression(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "regression.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	service := workflow.New(st)
	if _, err = service.Register(model.RecordInput{ID: "approved-2", StoreID: "main", Title: "Approved Choice", Score: 92}, 1, "clerk"); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Submit("approved-2", "clerk", 2); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Review("approved-2", "manager", true, "ready", 3); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Register(model.RecordInput{ID: "published-1", StoreID: "main", Title: "Published Choice", Score: 92}, 1, "clerk"); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Submit("published-1", "clerk", 2); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Review("published-1", "manager", true, "ready", 3); err != nil {
		t.Fatal(err)
	}
	if _, err = service.Publish("published-1", "manager", 4); err != nil {
		t.Fatal(err)
	}
	collection, err := service.Recommendations(recommendation.Request{StoreID: "main", IncludePublished: true, Limit: 1}, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(collection.Items) != 1 {
		t.Fatalf("expected one recommendation, got %d", len(collection.Items))
	}
	if collection.Items[0].Record.ID != "published-1" {
		t.Fatalf("expected published-1, got %s", collection.Items[0].Record.ID)
	}
}
