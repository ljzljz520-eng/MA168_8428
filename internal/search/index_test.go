package search

import (
	"bookstore/recommendation/internal/model"
	"testing"
)

func TestIndexSearch(t *testing.T) {
	index := New()
	index.Add(model.Record{ID: "a", Title: "Moon River", Summary: "quiet fiction"})
	index.Add(model.Record{ID: "b", Title: "Moon", Summary: "history"})
	hits := index.Search("moon", 1)
	if len(hits) != 1 || hits[0].ID != "a" {
		t.Fatalf("index result: %#v", hits)
	}
	index.Remove("a")
	if index.Count() != 1 {
		t.Fatal("remove failed")
	}
}
