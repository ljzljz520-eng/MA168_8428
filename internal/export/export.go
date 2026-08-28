package export

import (
	"bookstore/recommendation/internal/model"
	"encoding/json"
	"fmt"
	"io"
)

type Bundle struct {
	Records   []model.Record     `json:"records"`
	Audits    []model.AuditEvent `json:"audits"`
	Workflows []model.Workflow   `json:"workflows"`
}

func NewBundle(records []model.Record, audits []model.AuditEvent, workflows []model.Workflow) Bundle {
	return Bundle{Records: append([]model.Record(nil), records...), Audits: append([]model.AuditEvent(nil), audits...), Workflows: append([]model.Workflow(nil), workflows...)}
}

func WriteJSON(writer io.Writer, bundle Bundle) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(bundle)
}

func WriteCSV(writer io.Writer, records []model.Record) error {
	if _, err := io.WriteString(writer, "id,store_id,title,author,genre,score,status,version\n"); err != nil {
		return err
	}
	for _, record := range records {
		line := fmt.Sprintf("%s,%s,%s,%s,%s,%d,%s,%d\n", clean(record.ID), clean(record.StoreID), clean(record.Title), clean(record.Author), clean(record.Genre), record.Score, record.Status, record.Version)
		if _, err := io.WriteString(writer, line); err != nil {
			return err
		}
	}
	return nil
}

func clean(value string) string {
	result := ""
	for _, r := range value {
		if r == ',' || r == '\n' || r == '\r' {
			result += " "
		} else {
			result += string(r)
		}
	}
	return result
}
