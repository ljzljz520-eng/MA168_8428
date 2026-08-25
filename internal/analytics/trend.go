package analytics

import (
	"bookstore/recommendation/internal/model"
	"sort"
)

type Point struct {
	At      int64 `json:"at"`
	Count   int   `json:"count"`
	Average int   `json:"average"`
}

func Trend(records []model.Record) []Point {
	groups := map[int64][]model.Record{}
	for _, record := range records {
		groups[record.UpdatedAt] = append(groups[record.UpdatedAt], record)
	}
	points := make([]Point, 0, len(groups))
	for at, items := range groups {
		points = append(points, Point{At: at, Count: len(items), Average: ScoreAverage(items)})
	}
	sort.Slice(points, func(i, j int) bool { return points[i].At < points[j].At })
	return points
}

func Latest(records []model.Record) (model.Record, bool) {
	if len(records) == 0 {
		return model.Record{}, false
	}
	result := records[0]
	for _, record := range records[1:] {
		if record.UpdatedAt > result.UpdatedAt || record.UpdatedAt == result.UpdatedAt && record.ID > result.ID {
			result = record
		}
	}
	return result, true
}

func Growth(points []Point) int {
	if len(points) < 2 {
		return 0
	}
	first, last := points[0].Count, points[len(points)-1].Count
	if first == 0 {
		if last == 0 {
			return 0
		}
		return 100
	}
	return (last - first) * 100 / first
}

func Flatten(points []Point) []int {
	result := make([]int, 0, len(points))
	for _, point := range points {
		result = append(result, point.Count)
	}
	return result
}
