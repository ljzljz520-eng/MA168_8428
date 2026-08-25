package catalog

import (
	"bookstore/recommendation/internal/model"
	"bookstore/recommendation/internal/store"
	"path/filepath"
	"testing"
)

func TestRegisterSearchUpdateArchive(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "catalog.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	catalog := New(st)
	input := model.RecordInput{ID: "r1", StoreID: "s1", Title: "River", Author: "Writer", Genre: "Travel", Score: 70}
	if _, err = catalog.Register(input, 10); err != nil {
		t.Fatal(err)
	}
	page, err := catalog.Search(model.Query{Text: "river"})
	if err != nil || len(page.Items) != 1 {
		t.Fatalf("search failed: %v", err)
	}
	input.Score = 90
	if _, err = catalog.Update("r1", input, 11); err != nil {
		t.Fatal(err)
	}
	if _, err = catalog.Archive("r1", 12); err != nil {
		t.Fatal(err)
	}
	page, err = catalog.Search(model.Query{})
	if err != nil || len(page.Items) != 0 {
		t.Fatalf("archived record should be hidden")
	}
}
