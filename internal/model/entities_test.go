package model

import "testing"

func TestRecordCloneAndLabel(t *testing.T) {
	record := Record{ID: "r1", Title: "The Shelf", Author: "A. Writer", Tags: []string{"fiction"}}
	clone := record.Clone()
	clone.Tags[0] = "history"
	if record.Tags[0] != "fiction" || record.DisplayLabel() != "The Shelf - A. Writer" {
		t.Fatalf("unexpected clone or label")
	}
}

func TestQueryMatchesAndPagination(t *testing.T) {
	records := []Record{{ID: "a", StoreID: "s", Title: "Moon", Genre: "Fiction", Score: 90, Status: StatusApproved}, {ID: "b", StoreID: "s", Title: "Sun", Genre: "History", Score: 60, Status: StatusArchived}}
	query := Query{StoreID: "s", Text: "moon"}
	if !query.Matches(records[0]) || query.Matches(records[1]) {
		t.Fatalf("query matching failed")
	}
	page := Paginate(records, Query{Limit: 1})
	if page.Total != 2 || len(page.Items) != 1 || !page.HasMore {
		t.Fatalf("pagination failed")
	}
}
