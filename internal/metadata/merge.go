package metadata

import (
	"bookstore/recommendation/internal/model"
	"fmt"
	"strings"
)

type MergeResult struct {
	Record    model.Record
	Conflicts []string
}

func Merge(base, incoming model.Record) MergeResult {
	result := base.Clone()
	conflicts := make([]string, 0)
	if incoming.Title != "" && incoming.Title != base.Title {
		result.Title = incoming.Title
	}
	if incoming.Author != "" && incoming.Author != base.Author {
		result.Author = incoming.Author
	}
	if incoming.Genre != "" && incoming.Genre != base.Genre {
		result.Genre = incoming.Genre
	}
	if incoming.Summary != "" && incoming.Summary != base.Summary {
		result.Summary = incoming.Summary
	}
	if incoming.Score != 0 && incoming.Score != base.Score {
		result.Score = incoming.Score
	}
	if incoming.Status != "" && incoming.Status != base.Status {
		conflicts = append(conflicts, fmt.Sprintf("status:%s->%s", base.Status, incoming.Status))
	}
	result.Tags = append(result.Tags, incoming.Tags...)
	result = MergeTags(result)
	return MergeResult{Record: result, Conflicts: conflicts}
}

func Canonical(record model.Record) model.Record {
	record.Title = strings.TrimSpace(record.Title)
	record.Author = strings.TrimSpace(record.Author)
	record.Genre = strings.ToLower(strings.TrimSpace(record.Genre))
	record.Summary = strings.TrimSpace(record.Summary)
	record.Tags = MergeTags(record).Tags
	return record
}

func EqualContent(a, b model.Record) bool {
	a = Canonical(a)
	b = Canonical(b)
	return a.Title == b.Title && a.Author == b.Author && a.Genre == b.Genre && a.Summary == b.Summary && a.Score == b.Score && strings.Join(a.Tags, ",") == strings.Join(b.Tags, ",")
}

func EmptyProfile() Profile {
	return Profile{Language: "zh-CN", Audience: "general", Shelf: "unsorted", Keywords: []string{}}
}
