package recommendation

import (
	"bookstore/recommendation/internal/model"
	"fmt"
	"sort"
)

type Explanation struct {
	RecordID string   `json:"record_id"`
	Headline string   `json:"headline"`
	Details  []string `json:"details"`
}

func Explain(candidate Candidate) Explanation {
	details := append([]string(nil), candidate.Reasons...)
	sort.Strings(details)
	return Explanation{RecordID: candidate.Record.ID, Headline: fmt.Sprintf("%s (%d/100)", candidate.Record.DisplayLabel(), candidate.Record.Score), Details: details}
}

func ExplainAll(candidates []Candidate) []Explanation {
	result := make([]Explanation, 0, len(candidates))
	for _, candidate := range candidates {
		result = append(result, Explain(candidate))
	}
	return result
}

func Compare(a, b model.Record) int {
	if a.Score != b.Score {
		if a.Score > b.Score {
			return 1
		}
		return -1
	}
	if statusRank(a.Status) != statusRank(b.Status) {
		if statusRank(a.Status) > statusRank(b.Status) {
			return 1
		}
		return -1
	}
	if a.UpdatedAt != b.UpdatedAt {
		if a.UpdatedAt > b.UpdatedAt {
			return 1
		}
		return -1
	}
	if a.ID == b.ID {
		return 0
	}
	if a.ID < b.ID {
		return 1
	}
	return -1
}
