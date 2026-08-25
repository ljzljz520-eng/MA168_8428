package workflow

import (
	"bookstore/recommendation/internal/catalog"
	"bookstore/recommendation/internal/model"
	"bookstore/recommendation/internal/policy"
	"bookstore/recommendation/internal/recommendation"
	"bookstore/recommendation/internal/store"
	"fmt"
	"strings"
)

type Service struct {
	Catalog  *catalog.Catalog
	Store    *store.Store
	Policy   *policy.Policy
	sequence int
}

func New(st *store.Store) *Service {
	return &Service{Store: st, Catalog: catalog.New(st), Policy: policy.New()}
}

func (s *Service) nextID(prefix string) string {
	s.sequence++
	return fmt.Sprintf("%s-%03d", prefix, s.sequence)
}

func (s *Service) Register(input model.RecordInput, at int64, actor string) (model.Record, error) {
	record, err := s.Catalog.Register(input, at)
	if err != nil {
		return model.Record{}, err
	}
	event := model.NewAuditEvent(s.nextID("audit"), record.ID, actor, "register", at, "", record.Status, "record registered")
	if err = s.Store.SaveAudit(event); err != nil {
		return model.Record{}, err
	}
	if err = s.Store.SaveWorkflow(model.NewWorkflow(s.nextID("flow"), record.ID, "registration", "registered", 1, at)); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Submit(id, actor string, at int64) (model.Record, error) {
	return s.transition(id, model.StatusPending, "submit", actor, at, "submitted for review")
}

func (s *Service) Review(id, actor string, approve bool, note string, at int64) (model.Record, error) {
	target := model.StatusRejected
	action := "reject"
	if approve {
		target, action = model.StatusApproved, "approve"
	}
	record, err := s.transition(id, target, action, actor, at, note)
	if err != nil {
		return model.Record{}, err
	}
	record.Reviewer, record.ReviewNote = actor, note
	if strings.TrimSpace(record.ReviewNote) == "" {
		record.ReviewNote = s.Policy.Explain(record)
	}
	if err = s.Store.SaveRecord(record); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Publish(id, actor string, at int64) (model.Record, error) {
	return s.transition(id, model.StatusPublished, "publish", actor, at, "published to store")
}

func (s *Service) Archive(id, actor string, at int64) (model.Record, error) {
	record, err := s.Catalog.Archive(id, at)
	if err != nil {
		return model.Record{}, err
	}
	if err = s.Store.SaveAudit(model.NewAuditEvent(s.nextID("audit"), id, actor, "archive", at, record.Status, model.StatusArchived, "record archived")); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Update(id string, input model.RecordInput, actor string, at int64) (model.Record, error) {
	record, err := s.Catalog.Update(id, input, at)
	if err != nil {
		return model.Record{}, err
	}
	if err = s.Store.SaveAudit(model.NewAuditEvent(s.nextID("audit"), id, actor, "update", at, record.Status, record.Status, "details updated")); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) transition(id string, target model.Status, action, actor string, at int64, note string) (model.Record, error) {
	record, err := s.Store.GetRecord(id)
	if err != nil {
		return model.Record{}, err
	}
	from := record.Status
	if err = model.ValidateTransition(from, target); err != nil {
		return model.Record{}, err
	}
	record.Status, record.Version, record.UpdatedAt = target, record.Version+1, at
	event := model.NewAuditEvent(s.nextID("audit"), id, actor, action, at, from, target, note)
	if err = s.Store.SaveRecordWithAudit(record, event); err != nil {
		return model.Record{}, err
	}
	return record, nil
}

func (s *Service) Recommendations(request recommendation.Request, at int64) (recommendation.Collection, error) {
	records, err := s.Store.ListRecords()
	if err != nil {
		return recommendation.Collection{}, err
	}
	return recommendation.NewCollection(recommendation.Recommend(records, request), at), nil
}
