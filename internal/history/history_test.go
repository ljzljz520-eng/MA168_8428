package history

import (
	"bookstore/recommendation/internal/model"
	"testing"
)

func TestTimelineAndDiff(t *testing.T) {
	events := []model.AuditEvent{{RecordID: "r", Action: "register", Actor: "a", At: 1, To: model.StatusDraft}, {RecordID: "r", Action: "approve", Actor: "b", At: 3, To: model.StatusApproved}}
	timeline := Build("r", events)
	if timeline.LastStatus() != model.StatusApproved || !timeline.Contains("approve") || timeline.Duration() != 2 {
		t.Fatal("timeline failed")
	}
	changes := Diff(model.Record{Title: "old", Score: 1, Version: 1}, model.Record{Title: "new", Score: 2, Version: 2})
	if len(changes) != 3 {
		t.Fatalf("diff failed: %#v", changes)
	}
}
