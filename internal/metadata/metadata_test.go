package metadata

import (
	"bookstore/recommendation/internal/model"
	"testing"
)

func TestEnrichment(t *testing.T) {
	record := model.Record{ID: "r", StoreID: "s", Title: "Title", Author: "Author", Genre: "fiction", Summary: "Summary", Score: 95}
	enriched := Enrich(record)
	if enriched.Completeness != 100 || enriched.Profile.Audience != "featured" {
		t.Fatal("enrichment failed")
	}
	merged := MergeTags(record, "new", "new")
	if len(merged.Tags) != 1 {
		t.Fatal("tag merge failed")
	}
}
