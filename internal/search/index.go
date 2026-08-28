package search

import (
	"bookstore/recommendation/internal/model"
	"sort"
	"strings"
)

type Hit struct {
	ID    string   `json:"id"`
	Score int      `json:"score"`
	Terms []string `json:"terms"`
}

type Index struct {
	records map[string]model.Record
	terms   map[string]map[string]bool
}

func New() *Index {
	return &Index{records: map[string]model.Record{}, terms: map[string]map[string]bool{}}
}

func (i *Index) Add(record model.Record) {
	i.records[record.ID] = record
	for _, term := range tokenize(record.Title + " " + record.Author + " " + record.Summary) {
		if i.terms[term] == nil {
			i.terms[term] = map[string]bool{}
		}
		i.terms[term][record.ID] = true
	}
}

func (i *Index) Remove(id string) {
	delete(i.records, id)
	for term, ids := range i.terms {
		delete(ids, id)
		if len(ids) == 0 {
			delete(i.terms, term)
		}
	}
}

func (i *Index) Search(text string, limit int) []Hit {
	terms := tokenize(text)
	candidates := map[string]int{}
	for _, term := range terms {
		for id := range i.terms[term] {
			candidates[id]++
		}
	}
	result := make([]Hit, 0, len(candidates))
	for id, score := range candidates {
		result = append(result, Hit{ID: id, Score: score, Terms: append([]string(nil), terms...)})
	}
	sort.Slice(result, func(a, b int) bool {
		if result[a].Score == result[b].Score {
			return result[a].ID < result[b].ID
		}
		return result[a].Score > result[b].Score
	})
	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}
	return result
}

func tokenize(text string) []string {
	raw := strings.Fields(strings.ToLower(text))
	result := make([]string, 0, len(raw))
	for _, value := range raw {
		value = strings.Trim(value, ".,!?;:()[]{}\"")
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func (i *Index) Count() int { return len(i.records) }
