package workflow

import (
	"bookstore/recommendation/internal/catalog"
	"bookstore/recommendation/internal/history"
	"bookstore/recommendation/internal/model"
	"sort"
)

type Report struct {
	Summary  catalog.Summary  `json:"summary"`
	Facets   catalog.Facets   `json:"facets"`
	Audits   int              `json:"audits"`
	Timeline history.Timeline `json:"timeline"`
}

func (s *Service) BuildReport(recordID string) (Report, error) {
	records, err := s.Store.ListRecords()
	if err != nil {
		return Report{}, err
	}
	audits, err := s.Store.ListAudits(recordID)
	if err != nil {
		return Report{}, err
	}
	if recordID != "" {
		filtered := make([]model.Record, 0, 1)
		for _, record := range records {
			if record.ID == recordID {
				filtered = append(filtered, record)
			}
		}
		records = filtered
	}
	return Report{Summary: catalog.Summarize(records), Facets: catalog.BuildFacets(records), Audits: len(audits), Timeline: history.Build(recordID, audits)}, nil
}

func SortByUpdated(records []model.Record) []model.Record {
	copyRecords := append([]model.Record(nil), records...)
	sort.Slice(copyRecords, func(i, j int) bool { return copyRecords[i].UpdatedAt > copyRecords[j].UpdatedAt })
	return copyRecords
}
