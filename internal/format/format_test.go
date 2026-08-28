package format

import (
	"bookstore/recommendation/internal/model"
	"strings"
	"testing"
)

func TestFormatting(t *testing.T) {
	record := model.Record{ID: "r", Title: "Book", Author: "A", Score: 90, Status: model.StatusApproved, Tags: []string{" fiction "}}
	if !strings.Contains(RecordLine(record), "Book") || StatusLabel(record.Status) != "已批准" || TagsLine(record.Tags) != "fiction" {
		t.Fatal("format failed")
	}
}
