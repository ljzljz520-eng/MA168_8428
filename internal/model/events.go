package model

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
)

func NewAuditEvent(id, recordID, actor, action string, at int64, from, to Status, note string) AuditEvent {
	return AuditEvent{ID: id, RecordID: recordID, Actor: actor, Action: action, At: at, From: from, To: to, Note: note}
}

func NewWorkflow(id, recordID, name, step string, position int, at int64) Workflow {
	return Workflow{ID: id, RecordID: recordID, Name: name, Step: step, Position: position, At: at}
}

func NewAttachment(id, recordID, name, mediaType string, content []byte) Attachment {
	h := sha256.Sum256(content)
	return Attachment{ID: id, RecordID: recordID, Name: name, MediaType: mediaType, Content: append([]byte(nil), content...), Checksum: hex.EncodeToString(h[:])}
}

func SortTags(tags []string) []string {
	copyTags := append([]string(nil), tags...)
	sort.Strings(copyTags)
	return copyTags
}
