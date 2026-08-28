package metadata

import (
	"bookstore/recommendation/internal/model"
	"sort"
	"strings"
)

type Profile struct {
	Language string   `json:"language"`
	Audience string   `json:"audience"`
	Shelf    string   `json:"shelf"`
	Keywords []string `json:"keywords"`
}

type Enriched struct {
	Record       model.Record `json:"record"`
	Profile      Profile      `json:"profile"`
	Completeness int          `json:"completeness"`
}

func Infer(record model.Record) Profile {
	profile := Profile{Language: "zh-CN", Audience: "general", Shelf: strings.ToLower(record.Genre), Keywords: []string{strings.ToLower(record.Title), strings.ToLower(record.Author)}}
	if record.Score >= 90 {
		profile.Audience = "featured"
	}
	if record.Status == model.StatusArchived {
		profile.Shelf = "archive"
	}
	return profile
}

func Enrich(record model.Record) Enriched {
	profile := Infer(record)
	completeness := 0
	fields := []string{record.ID, record.StoreID, record.Title, record.Author, record.Genre, record.Summary}
	for _, field := range fields {
		if strings.TrimSpace(field) != "" {
			completeness += 1
		}
	}
	return Enriched{Record: record, Profile: profile, Completeness: completeness * 100 / len(fields)}
}

func Keywords(records []model.Record) []string {
	set := map[string]bool{}
	for _, record := range records {
		for _, value := range strings.Fields(strings.ToLower(record.Title + " " + record.Author + " " + record.Genre)) {
			set[value] = true
		}
	}
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func MergeTags(record model.Record, extra ...string) model.Record {
	record.Tags = append([]string(nil), record.Tags...)
	for _, tag := range extra {
		tag = strings.TrimSpace(tag)
		if tag != "" && !contains(record.Tags, tag) {
			record.Tags = append(record.Tags, tag)
		}
	}
	sort.Strings(record.Tags)
	return record
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(value, target) {
			return true
		}
	}
	return false
}
