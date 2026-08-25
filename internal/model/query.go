package model

import "strings"

type Query struct {
	StoreID         string
	Text            string
	Genre           string
	Status          Status
	MinScore        int
	MaxScore        int
	IncludeArchived bool
	Offset          int
	Limit           int
}

func (q Query) Normalize() Query {
	q.StoreID = strings.TrimSpace(q.StoreID)
	q.Text = strings.ToLower(strings.TrimSpace(q.Text))
	q.Genre = strings.ToLower(strings.TrimSpace(q.Genre))
	if q.Offset < 0 {
		q.Offset = 0
	}
	if q.Limit <= 0 || q.Limit > 100 {
		q.Limit = 20
	}
	if q.MaxScore == 0 {
		q.MaxScore = 100
	}
	return q
}

func (q Query) Matches(r Record) bool {
	q = q.Normalize()
	if q.StoreID != "" && r.StoreID != q.StoreID {
		return false
	}
	if q.Status != "" && r.Status != q.Status {
		return false
	}
	if !q.IncludeArchived && r.Status == StatusArchived {
		return false
	}
	if r.Score < q.MinScore || r.Score > q.MaxScore {
		return false
	}
	if q.Genre != "" && strings.ToLower(r.Genre) != q.Genre {
		return false
	}
	if q.Text != "" {
		text := strings.ToLower(r.Title + " " + r.Author + " " + r.Summary)
		if !strings.Contains(text, q.Text) {
			return false
		}
	}
	return true
}

type Page struct {
	Items   []Record `json:"items"`
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
	Total   int      `json:"total"`
	HasMore bool     `json:"has_more"`
}

func Paginate(records []Record, q Query) Page {
	q = q.Normalize()
	start := q.Offset
	if start > len(records) {
		start = len(records)
	}
	end := start + q.Limit
	if end > len(records) {
		end = len(records)
	}
	items := make([]Record, end-start)
	copy(items, records[start:end])
	return Page{Items: items, Offset: start, Limit: q.Limit, Total: len(records), HasMore: end < len(records)}
}
