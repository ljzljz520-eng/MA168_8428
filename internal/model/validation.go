package model

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrMissingID     = errors.New("record id is required")
	ErrMissingStore  = errors.New("store id is required")
	ErrMissingTitle  = errors.New("title is required")
	ErrInvalidScore  = errors.New("score must be between 0 and 100")
	ErrInvalidStatus = errors.New("invalid record status")
)

func ValidateInput(input RecordInput) error {
	input = input.Normalize()
	if input.ID == "" {
		return ErrMissingID
	}
	if input.StoreID == "" {
		return ErrMissingStore
	}
	if input.Title == "" {
		return ErrMissingTitle
	}
	if input.Score < 0 || input.Score > 100 {
		return ErrInvalidScore
	}
	if len(input.Title) > 200 {
		return fmt.Errorf("title exceeds 200 characters")
	}
	if len(input.Summary) > 2000 {
		return fmt.Errorf("summary exceeds 2000 characters")
	}
	return nil
}

func ValidateTransition(from, to Status) error {
	if from == to {
		return errors.New("status is unchanged")
	}
	allowed := map[Status][]Status{
		StatusDraft:     {StatusPending, StatusArchived},
		StatusPending:   {StatusApproved, StatusRejected, StatusDraft},
		StatusApproved:  {StatusPublished, StatusDraft, StatusArchived},
		StatusPublished: {StatusArchived, StatusDraft},
		StatusRejected:  {StatusDraft, StatusArchived},
		StatusArchived:  {},
	}
	for _, candidate := range allowed[from] {
		if candidate == to {
			return nil
		}
	}
	return fmt.Errorf("cannot transition from %s to %s", from, to)
}

func ValidateStatus(status Status) error {
	if strings.TrimSpace(string(status)) == "" {
		return ErrInvalidStatus
	}
	for _, value := range []Status{StatusDraft, StatusPending, StatusApproved, StatusRejected, StatusPublished, StatusArchived} {
		if status == value {
			return nil
		}
	}
	return ErrInvalidStatus
}
