package format

import (
	"bookstore/recommendation/internal/model"
	"fmt"
	"strings"
)

func RecordLine(record model.Record) string {
	return fmt.Sprintf("%s | %s | %s | %d | %s", record.ID, record.Title, record.Author, record.Score, record.Status)
}

func Table(records []model.Record) string {
	lines := []string{"ID | TITLE | AUTHOR | SCORE | STATUS"}
	for _, record := range records {
		lines = append(lines, RecordLine(record))
	}
	return strings.Join(lines, "\n")
}

func StatusLabel(status model.Status) string {
	switch status {
	case model.StatusDraft:
		return "草稿"
	case model.StatusPending:
		return "待审核"
	case model.StatusApproved:
		return "已批准"
	case model.StatusRejected:
		return "已拒绝"
	case model.StatusPublished:
		return "已发布"
	case model.StatusArchived:
		return "已归档"
	default:
		return "未知"
	}
}

func ScoreLabel(score int) string {
	switch {
	case score >= 90:
		return "重点推荐"
	case score >= 75:
		return "推荐"
	case score >= 60:
		return "可观察"
	default:
		return "待补充"
	}
}

func TagsLine(tags []string) string {
	cleaned := make([]string, 0, len(tags))
	for _, tag := range tags {
		if value := strings.TrimSpace(tag); value != "" {
			cleaned = append(cleaned, value)
		}
	}
	return strings.Join(cleaned, ", ")
}

func Summaries(records []model.Record) []string {
	result := make([]string, 0, len(records))
	for _, record := range records {
		result = append(result, record.Title+": "+ScoreLabel(record.Score)+" / "+StatusLabel(record.Status))
	}
	return result
}

func Header(storeID string, at int64) string {
	return fmt.Sprintf("store=%s generated=%d", strings.TrimSpace(storeID), at)
}

func JoinLines(lines []string) string {
	cleaned := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	return strings.Join(cleaned, "\n")
}
