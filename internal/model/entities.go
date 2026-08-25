package model

import "strings"

type Status string

const (
	StatusDraft     Status = "draft"
	StatusPending   Status = "pending"
	StatusApproved  Status = "approved"
	StatusRejected  Status = "rejected"
	StatusPublished Status = "published"
	StatusArchived  Status = "archived"
)

type Record struct {
	ID         string   `json:"id"`
	StoreID    string   `json:"store_id"`
	Title      string   `json:"title"`
	Author     string   `json:"author"`
	Genre      string   `json:"genre"`
	Summary    string   `json:"summary"`
	Score      int      `json:"score"`
	Status     Status   `json:"status"`
	Version    int      `json:"version"`
	CreatedAt  int64    `json:"created_at"`
	UpdatedAt  int64    `json:"updated_at"`
	Reviewer   string   `json:"reviewer"`
	ReviewNote string   `json:"review_note"`
	Tags       []string `json:"tags"`
}

type AuditEvent struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	At       int64  `json:"at"`
	From     Status `json:"from"`
	To       Status `json:"to"`
	Note     string `json:"note"`
}

type Workflow struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Name      string `json:"name"`
	Step      string `json:"step"`
	Position  int    `json:"position"`
	Completed bool   `json:"completed"`
	At        int64  `json:"at"`
}

type Attachment struct {
	ID        string `json:"id"`
	RecordID  string `json:"record_id"`
	Name      string `json:"name"`
	MediaType string `json:"media_type"`
	Content   []byte `json:"content"`
	Checksum  string `json:"checksum"`
}

type RecordInput struct {
	ID      string   `json:"id"`
	StoreID string   `json:"store_id"`
	Title   string   `json:"title"`
	Author  string   `json:"author"`
	Genre   string   `json:"genre"`
	Summary string   `json:"summary"`
	Score   int      `json:"score"`
	Tags    []string `json:"tags"`
}

func (r RecordInput) Normalize() RecordInput {
	r.ID = strings.TrimSpace(r.ID)
	r.StoreID = strings.TrimSpace(r.StoreID)
	r.Title = strings.TrimSpace(r.Title)
	r.Author = strings.TrimSpace(r.Author)
	r.Genre = strings.TrimSpace(r.Genre)
	r.Summary = strings.TrimSpace(r.Summary)
	for i := range r.Tags {
		r.Tags[i] = strings.TrimSpace(r.Tags[i])
	}
	return r
}

func (r Record) Clone() Record {
	r.Tags = append([]string(nil), r.Tags...)
	return r
}

func (r Record) IsVisible() bool { return r.Status == StatusApproved || r.Status == StatusPublished }

func (r Record) IsTerminal() bool { return r.Status == StatusArchived || r.Status == StatusRejected }

func (r Record) DisplayLabel() string {
	if r.Author == "" {
		return r.Title
	}
	return r.Title + " - " + r.Author
}
