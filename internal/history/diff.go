package history

import (
	"bookstore/recommendation/internal/model"
	"fmt"
)

type Change struct {
	Field  string `json:"field"`
	Before string `json:"before"`
	After  string `json:"after"`
}

func Diff(before, after model.Record) []Change {
	changes := make([]Change, 0)
	if before.Title != after.Title {
		changes = append(changes, Change{"title", before.Title, after.Title})
	}
	if before.Author != after.Author {
		changes = append(changes, Change{"author", before.Author, after.Author})
	}
	if before.Genre != after.Genre {
		changes = append(changes, Change{"genre", before.Genre, after.Genre})
	}
	if before.Summary != after.Summary {
		changes = append(changes, Change{"summary", before.Summary, after.Summary})
	}
	if before.Score != after.Score {
		changes = append(changes, Change{"score", fmt.Sprint(before.Score), fmt.Sprint(after.Score)})
	}
	if before.Status != after.Status {
		changes = append(changes, Change{"status", string(before.Status), string(after.Status)})
	}
	if before.Version != after.Version {
		changes = append(changes, Change{"version", fmt.Sprint(before.Version), fmt.Sprint(after.Version)})
	}
	return changes
}

func Describe(changes []Change) string {
	if len(changes) == 0 {
		return "no changes"
	}
	result := ""
	for i, change := range changes {
		if i > 0 {
			result += "; "
		}
		result += change.Field + "=" + change.After
	}
	return result
}
