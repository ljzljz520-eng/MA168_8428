package model

import "testing"

func TestValidationAndTransitions(t *testing.T) {
	if ValidateInput(RecordInput{ID: "", StoreID: "s", Title: "x", Score: 10}) == nil {
		t.Fatal("missing id accepted")
	}
	if ValidateInput(RecordInput{ID: "a", StoreID: "s", Title: "x", Score: 101}) == nil {
		t.Fatal("invalid score accepted")
	}
	if ValidateTransition(StatusDraft, StatusPending) != nil {
		t.Fatal("draft should submit")
	}
	if ValidateTransition(StatusArchived, StatusDraft) == nil {
		t.Fatal("archived should be terminal")
	}
}
