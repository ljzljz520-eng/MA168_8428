package catalog

import (
	"bookstore/recommendation/internal/model"
	"bookstore/recommendation/internal/store"
	"fmt"
	"sort"
)

type Catalog struct {
	store *store.Store
}

func New(st *store.Store) *Catalog { return &Catalog{store: st} }

func (c *Catalog) Register(input model.RecordInput, at int64) (model.Record, error) {
	input = input.Normalize()
	if err := model.ValidateInput(input); err != nil {
		return model.Record{}, err
	}
	if _, err := c.store.GetRecord(input.ID); err == nil {
		return model.Record{}, fmt.Errorf("record %s already exists", input.ID)
	}
	record := model.Record{ID: input.ID, StoreID: input.StoreID, Title: input.Title, Author: input.Author, Genre: input.Genre, Summary: input.Summary, Score: input.Score, Status: model.StatusDraft, Version: 1, CreatedAt: at, UpdatedAt: at, Tags: model.SortTags(input.Tags)}
	if err := c.store.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (c *Catalog) Update(id string, input model.RecordInput, at int64) (model.Record, error) {
	record, err := c.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	input.ID = id
	input = input.Normalize()
	if err := model.ValidateInput(input); err != nil {
		return model.Record{}, err
	}
	if record.IsTerminal() {
		return model.Record{}, fmt.Errorf("terminal record cannot be changed")
	}
	record.StoreID, record.Title, record.Author, record.Genre = input.StoreID, input.Title, input.Author, input.Genre
	record.Summary, record.Score, record.Tags = input.Summary, input.Score, model.SortTags(input.Tags)
	record.Version++
	record.UpdatedAt = at
	if err := c.store.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (c *Catalog) Get(id string) (model.Record, error) { return c.store.GetRecord(id) }

func (c *Catalog) Search(query model.Query) (model.Page, error) {
	items, err := c.store.ListRecords()
	if err != nil {
		return model.Page{}, err
	}
	filtered := make([]model.Record, 0, len(items))
	for _, item := range items {
		if query.Matches(item) {
			filtered = append(filtered, item)
		}
	}
	sort.SliceStable(filtered, func(i, j int) bool {
		if filtered[i].Score == filtered[j].Score {
			return filtered[i].Title < filtered[j].Title
		}
		return filtered[i].Score > filtered[j].Score
	})
	return model.Paginate(filtered, query), nil
}

func (c *Catalog) Archive(id string, at int64) (model.Record, error) {
	record, err := c.store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	if err := model.ValidateTransition(record.Status, model.StatusArchived); err != nil {
		return model.Record{}, err
	}
	record.Status = model.StatusArchived
	record.Version++
	record.UpdatedAt = at
	if err := c.store.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}
