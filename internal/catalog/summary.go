package catalog

import (
	"bookstore/recommendation/internal/model"
	"fmt"
)

type Summary struct {
	Total        int `json:"total"`
	Visible      int `json:"visible"`
	Pending      int `json:"pending"`
	Archived     int `json:"archived"`
	AverageScore int `json:"average_score"`
}

func Summarize(records []model.Record) Summary {
	result := Summary{Total: len(records)}
	totalScore := 0
	for _, record := range records {
		totalScore += record.Score
		if record.IsVisible() {
			result.Visible++
		}
		if record.Status == model.StatusPending {
			result.Pending++
		}
		if record.Status == model.StatusArchived {
			result.Archived++
		}
	}
	if result.Total > 0 {
		result.AverageScore = totalScore / result.Total
	}
	return result
}

func FormatSummary(summary Summary) string {
	return fmt.Sprintf("total=%d visible=%d pending=%d archived=%d average=%d", summary.Total, summary.Visible, summary.Pending, summary.Archived, summary.AverageScore)
}
