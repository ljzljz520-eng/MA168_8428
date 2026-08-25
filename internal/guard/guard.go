package guard

import (
	"bookstore/recommendation/internal/model"
	"errors"
	"fmt"
	"strings"
)

var ErrUnsafe = errors.New("unsafe business value")

type Limits struct {
	MaxTitle      int
	MaxSummary    int
	MaxTags       int
	MaxAttachment int
}

func DefaultLimits() Limits {
	return Limits{MaxTitle: 200, MaxSummary: 2000, MaxTags: 20, MaxAttachment: 2 << 20}
}

func CheckRecord(record model.Record, limits Limits) error {
	if limits.MaxTitle <= 0 || limits.MaxSummary <= 0 {
		return fmt.Errorf("invalid limits")
	}
	if len(record.Title) > limits.MaxTitle || len(record.Summary) > limits.MaxSummary {
		return ErrUnsafe
	}
	if len(record.Tags) > limits.MaxTags {
		return ErrUnsafe
	}
	if record.Score < 0 || record.Score > 100 {
		return ErrUnsafe
	}
	return nil
}

func CheckAttachment(attachment model.Attachment, limits Limits) error {
	if attachment.ID == "" || attachment.RecordID == "" {
		return ErrUnsafe
	}
	if len(attachment.Content) > limits.MaxAttachment {
		return ErrUnsafe
	}
	if strings.TrimSpace(attachment.MediaType) == "" {
		return ErrUnsafe
	}
	return nil
}

func SanitizeText(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "\x00", "")
	value = strings.ReplaceAll(value, "\r\n", "\n")
	return value
}

func SanitizeInput(input model.RecordInput) model.RecordInput {
	input.ID = SanitizeText(input.ID)
	input.StoreID = SanitizeText(input.StoreID)
	input.Title = SanitizeText(input.Title)
	input.Author = SanitizeText(input.Author)
	input.Genre = SanitizeText(input.Genre)
	input.Summary = SanitizeText(input.Summary)
	for i := range input.Tags {
		input.Tags[i] = SanitizeText(input.Tags[i])
	}
	return input
}

func SafeID(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" {
		return false
	}
	for _, r := range value {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '_' {
			return false
		}
	}
	return true
}

func NormalizeScore(score int) int {
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}

func IsPrintable(value string) bool {
	for _, r := range value {
		if r < 32 && r != '\n' && r != '\t' {
			return false
		}
	}
	return true
}
