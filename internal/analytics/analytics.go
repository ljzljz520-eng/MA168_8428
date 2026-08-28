package analytics

import (
	"bookstore/recommendation/internal/model"
	"sort"
)

type Bucket struct {
	Label   string `json:"label"`
	Count   int    `json:"count"`
	Average int    `json:"average"`
}

type Dashboard struct {
	Total        int                  `json:"total"`
	ByStatus     map[model.Status]int `json:"by_status"`
	ByGenre      map[string]int       `json:"by_genre"`
	ScoreBuckets []Bucket             `json:"score_buckets"`
	TopRecords   []string             `json:"top_records"`
}

func Build(records []model.Record) Dashboard {
	dashboard := Dashboard{Total: len(records), ByStatus: map[model.Status]int{}, ByGenre: map[string]int{}}
	bucketScores := map[string][]int{"low": {}, "medium": {}, "high": {}}
	for _, record := range records {
		dashboard.ByStatus[record.Status]++
		dashboard.ByGenre[record.Genre]++
		switch {
		case record.Score < 50:
			bucketScores["low"] = append(bucketScores["low"], record.Score)
		case record.Score < 80:
			bucketScores["medium"] = append(bucketScores["medium"], record.Score)
		default:
			bucketScores["high"] = append(bucketScores["high"], record.Score)
		}
	}
	for _, label := range []string{"low", "medium", "high"} {
		scores := bucketScores[label]
		total := 0
		for _, score := range scores {
			total += score
		}
		average := 0
		if len(scores) > 0 {
			average = total / len(scores)
		}
		dashboard.ScoreBuckets = append(dashboard.ScoreBuckets, Bucket{Label: label, Count: len(scores), Average: average})
	}
	dashboard.TopRecords = topIDs(records, 5)
	return dashboard
}

func topIDs(records []model.Record, limit int) []string {
	copyRecords := append([]model.Record(nil), records...)
	sort.SliceStable(copyRecords, func(i, j int) bool {
		if copyRecords[i].Score == copyRecords[j].Score {
			return copyRecords[i].ID < copyRecords[j].ID
		}
		return copyRecords[i].Score > copyRecords[j].Score
	})
	if limit > 0 && len(copyRecords) > limit {
		copyRecords = copyRecords[:limit]
	}
	result := make([]string, 0, len(copyRecords))
	for _, record := range copyRecords {
		result = append(result, record.ID)
	}
	return result
}

func StatusCounts(records []model.Record) map[model.Status]int {
	result := map[model.Status]int{}
	for _, record := range records {
		result[record.Status]++
	}
	return result
}

func ScoreAverage(records []model.Record) int {
	if len(records) == 0 {
		return 0
	}
	total := 0
	for _, record := range records {
		total += record.Score
	}
	return total / len(records)
}
