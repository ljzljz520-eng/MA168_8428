package guard

import (
	"bookstore/recommendation/internal/model"
	"testing"
)

func TestGuardLimits(t *testing.T) {
	limits := DefaultLimits()
	if CheckRecord(model.Record{Title: "ok", Summary: "ok", Score: 10}, limits) != nil {
		t.Fatal("valid record rejected")
	}
	if CheckRecord(model.Record{Title: "bad", Score: 101}, limits) == nil {
		t.Fatal("invalid record accepted")
	}
	if SanitizeText(" x\x00 ") != "x" {
		t.Fatal("sanitize failed")
	}
}
