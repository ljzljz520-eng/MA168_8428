package policy

import "bookstore/recommendation/internal/model"

type Check struct {
	Name   string
	Passed bool
	Detail string
}

func Checks(record model.Record) []Check {
	return []Check{{Name: "identity", Passed: record.ID != "" && record.StoreID != "", Detail: "record and store identifiers"}, {Name: "content", Passed: record.Title != "" && record.Summary != "", Detail: "title and summary"}, {Name: "score", Passed: record.Score >= 0 && record.Score <= 100, Detail: "score range"}, {Name: "state", Passed: model.ValidateStatus(record.Status) == nil, Detail: "known workflow state"}}
}

func AllPassed(checks []Check) bool {
	for _, check := range checks {
		if !check.Passed {
			return false
		}
	}
	return true
}

func Failed(checks []Check) []Check {
	result := make([]Check, 0)
	for _, check := range checks {
		if !check.Passed {
			result = append(result, check)
		}
	}
	return result
}
