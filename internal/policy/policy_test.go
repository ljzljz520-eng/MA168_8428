package policy

import (
	"bookstore/recommendation/internal/model"
	"testing"
)

func TestPolicyDecisions(t *testing.T) {
	policy := New()
	good := model.Record{ID: "r", StoreID: "s", Title: "Title", Author: "Author", Genre: "fiction", Score: 80}
	if !policy.Evaluate(good).Allowed {
		t.Fatal("good record rejected")
	}
	bad := good
	bad.Score = 10
	if policy.Evaluate(bad).Allowed {
		t.Fatal("low score accepted")
	}
	if !Can(Actor{ID: "m", Role: RoleManager}, "archive") || Can(Actor{ID: "c", Role: RoleClerk}, "archive") {
		t.Fatal("role matrix failed")
	}
}
