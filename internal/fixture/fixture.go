package fixture

import (
	"bookstore/recommendation/internal/model"
	"fmt"
	"strings"
)

type BookSpec struct {
	ID      string
	StoreID string
	Title   string
	Author  string
	Genre   string
	Score   int
	Summary string
	Status  model.Status
}

func Standard() []BookSpec {
	return []BookSpec{{"fiction-01", "central", "The Quiet Shelf", "M. Lin", "fiction", 88, "A patient story about a small shop.", model.StatusApproved}, {"history-01", "central", "Streets Remember", "J. Hu", "history", 82, "Local history told through shopkeepers.", model.StatusPublished}, {"essay-01", "east", "Notes on Reading", "S. Rao", "essay", 74, "Short essays for staff discussion.", model.StatusDraft}, {"poetry-01", "west", "Window Poems", "L. Chen", "poetry", 91, "Poems selected for the front table.", model.StatusApproved}}
}

func Inputs(specs []BookSpec) []model.RecordInput {
	result := make([]model.RecordInput, 0, len(specs))
	for _, spec := range specs {
		result = append(result, model.RecordInput{ID: spec.ID, StoreID: spec.StoreID, Title: spec.Title, Author: spec.Author, Genre: spec.Genre, Score: spec.Score, Summary: spec.Summary})
	}
	return result
}

func CSV(specs []BookSpec) string {
	lines := []string{"id,store,title,author,genre,score,summary"}
	for _, spec := range specs {
		lines = append(lines, fmt.Sprintf("%s,%s,%s,%s,%s,%d,%s", spec.ID, spec.StoreID, spec.Title, spec.Author, spec.Genre, spec.Score, spec.Summary))
	}
	return strings.Join(lines, "\n") + "\n"
}

func Lookup(specs []BookSpec, id string) (BookSpec, bool) {
	for _, spec := range specs {
		if spec.ID == id {
			return spec, true
		}
	}
	return BookSpec{}, false
}

func IDs(specs []BookSpec) []string {
	result := make([]string, 0, len(specs))
	for _, spec := range specs {
		result = append(result, spec.ID)
	}
	return result
}
