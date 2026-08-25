package history

import (
	"bookstore/recommendation/internal/model"
	"sort"
)

type Entry struct {
	At     int64  `json:"at"`
	Action string `json:"action"`
	Actor  string `json:"actor"`
	Status string `json:"status"`
	Note   string `json:"note"`
}

type Timeline struct {
	RecordID string  `json:"record_id"`
	Entries  []Entry `json:"entries"`
}

func Build(recordID string, events []model.AuditEvent) Timeline {
	entries := make([]Entry, 0, len(events))
	for _, event := range events {
		if recordID == "" || event.RecordID == recordID {
			entries = append(entries, Entry{At: event.At, Action: event.Action, Actor: event.Actor, Status: string(event.To), Note: event.Note})
		}
	}
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].At == entries[j].At {
			return entries[i].Action < entries[j].Action
		}
		return entries[i].At < entries[j].At
	})
	return Timeline{RecordID: recordID, Entries: entries}
}

func (t Timeline) LastStatus() model.Status {
	if len(t.Entries) == 0 {
		return ""
	}
	return model.Status(t.Entries[len(t.Entries)-1].Status)
}

func (t Timeline) Contains(action string) bool {
	for _, entry := range t.Entries {
		if entry.Action == action {
			return true
		}
	}
	return false
}

func (t Timeline) Duration() int64 {
	if len(t.Entries) < 2 {
		return 0
	}
	return t.Entries[len(t.Entries)-1].At - t.Entries[0].At
}
