package catalog

import (
	"bookstore/recommendation/internal/model"
	"strings"
)

type Facets struct {
	Genres       []string       `json:"genres"`
	Statuses     []model.Status `json:"statuses"`
	Stores       []string       `json:"stores"`
	ScoreBuckets map[string]int `json:"score_buckets"`
}

func BuildFacets(records []model.Record) Facets {
	genres, statuses, stores := map[string]bool{}, map[model.Status]bool{}, map[string]bool{}
	buckets := map[string]int{"0-49": 0, "50-79": 0, "80-100": 0}
	for _, record := range records {
		if record.Genre != "" {
			genres[strings.ToLower(record.Genre)] = true
		}
		statuses[record.Status] = true
		stores[record.StoreID] = true
		switch {
		case record.Score < 50:
			buckets["0-49"]++
		case record.Score < 80:
			buckets["50-79"]++
		default:
			buckets["80-100"]++
		}
	}
	result := Facets{ScoreBuckets: buckets}
	for genre := range genres {
		result.Genres = append(result.Genres, genre)
	}
	for status := range statuses {
		result.Statuses = append(result.Statuses, status)
	}
	for store := range stores {
		result.Stores = append(result.Stores, store)
	}
	return result
}

func MatchAll(records []model.Record, query model.Query) []model.Record {
	result := make([]model.Record, 0)
	for _, record := range records {
		if query.Matches(record) {
			result = append(result, record)
		}
	}
	return result
}
