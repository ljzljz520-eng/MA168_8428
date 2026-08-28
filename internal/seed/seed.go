package seed

import (
	"bookstore/recommendation/internal/fixture"
	"bookstore/recommendation/internal/model"
	"bookstore/recommendation/internal/workflow"
	"fmt"
)

type Result struct {
	IDs    []string
	Errors []string
}

type Loader struct{ service *workflow.Service }

func New(service *workflow.Service) *Loader { return &Loader{service: service} }

func (l *Loader) Load(specs []fixture.BookSpec, actor string, at int64) Result {
	result := Result{IDs: []string{}, Errors: []string{}}
	for index, spec := range specs {
		record, err := l.service.Register(model.RecordInput{ID: spec.ID, StoreID: spec.StoreID, Title: spec.Title, Author: spec.Author, Genre: spec.Genre, Score: spec.Score, Summary: spec.Summary}, at+int64(index), actor)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %v", spec.ID, err))
			continue
		}
		result.IDs = append(result.IDs, record.ID)
		if spec.Status == model.StatusPending || spec.Status == model.StatusApproved || spec.Status == model.StatusPublished {
			if _, err = l.service.Submit(record.ID, actor, at+int64(index)+1); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("%s submit: %v", spec.ID, err))
				continue
			}
		}
		if spec.Status == model.StatusApproved || spec.Status == model.StatusPublished {
			if _, err = l.service.Review(record.ID, actor, true, "seed approval", at+int64(index)+2); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("%s review: %v", spec.ID, err))
				continue
			}
		}
		if spec.Status == model.StatusPublished {
			if _, err = l.service.Publish(record.ID, actor, at+int64(index)+3); err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("%s publish: %v", spec.ID, err))
			}
		}
	}
	return result
}

func (l *Loader) LoadStandard(actor string, at int64) Result {
	return l.Load(fixture.Standard(), actor, at)
}

func Validate(result Result, expected int) error {
	if len(result.Errors) > 0 {
		return fmt.Errorf("seed errors: %v", result.Errors)
	}
	if len(result.IDs) != expected {
		return fmt.Errorf("seeded %d records, expected %d", len(result.IDs), expected)
	}
	return nil
}
