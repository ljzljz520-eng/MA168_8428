package history

import (
	"bookstore/recommendation/internal/model"
	"strings"
)

func Filter(events []model.AuditEvent, actor, action string) []model.AuditEvent {
	result := make([]model.AuditEvent, 0)
	for _, event := range events {
		if actor != "" && !strings.EqualFold(actor, event.Actor) {
			continue
		}
		if action != "" && event.Action != action {
			continue
		}
		result = append(result, event)
	}
	return result
}

func Actions(events []model.AuditEvent) []string {
	result := make([]string, 0, len(events))
	for _, event := range events {
		result = append(result, event.Action)
	}
	return result
}
