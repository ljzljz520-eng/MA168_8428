package recommendation

import (
	"bookstore/recommendation/internal/model"
	"testing"
)

func TestRecommendationFiltersAndRanks(t *testing.T) {
	records := []model.Record{{ID: "a", StoreID: "s", Genre: "fiction", Score: 80, Status: model.StatusApproved}, {ID: "b", StoreID: "s", Genre: "history", Score: 95, Status: model.StatusApproved}, {ID: "c", StoreID: "other", Genre: "fiction", Score: 99, Status: model.StatusPublished}}
	items := Recommend(records, Request{StoreID: "s", Genre: "fiction", Limit: 5})
	if len(items) != 1 || items[0].Record.ID != "a" {
		t.Fatalf("filter failed: %#v", items)
	}
	if Compare(records[1], records[0]) <= 0 {
		t.Fatal("comparison should favor higher score")
	}
}
