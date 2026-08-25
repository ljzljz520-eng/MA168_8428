package recommendation

import (
	"bookstore/recommendation/internal/model"
	"sort"
	"strings"
)

type Candidate struct {
	Record  model.Record `json:"record"`
	Reasons []string     `json:"reasons"`
	Rank    int          `json:"rank"`
}

type Request struct {
	StoreID          string
	Genre            string
	Limit            int
	IncludePublished bool
}

func NormalizeRequest(request Request) Request {
	request.StoreID = strings.TrimSpace(request.StoreID)
	request.Genre = strings.ToLower(strings.TrimSpace(request.Genre))
	if request.Limit <= 0 {
		request.Limit = 5
	}
	if request.Limit > 50 {
		request.Limit = 50
	}
	return request
}

func BuildCandidates(records []model.Record, request Request) []Candidate {
	request = NormalizeRequest(request)
	result := make([]Candidate, 0, len(records))
	for _, record := range records {
		if request.StoreID != "" && record.StoreID != request.StoreID {
			continue
		}
		if request.Genre != "" && strings.ToLower(record.Genre) != request.Genre {
			continue
		}
		if record.Status != model.StatusApproved && record.Status != model.StatusPublished {
			continue
		}
		if !request.IncludePublished && record.Status == model.StatusPublished {
			continue
		}
		reasons := []string{"score match"}
		if record.Status == model.StatusPublished {
			reasons = append(reasons, "published")
		}
		if record.Genre != "" {
			reasons = append(reasons, "genre:"+record.Genre)
		}
		result = append(result, Candidate{Record: record, Reasons: reasons})
	}
	return result
}

func statusRank(status model.Status) int {
	switch status {
	case model.StatusPublished:
		return 2
	case model.StatusApproved:
		return 1
	default:
		return 0
	}
}

func rankCandidate(candidate Candidate) int {
	return candidate.Record.Score*10 + statusRank(candidate.Record.Status)
}

func sortCandidates(candidates []Candidate) {
	sort.SliceStable(candidates, func(i, j int) bool {
		left, right := candidates[i].Record, candidates[j].Record
		if left.Score != right.Score {
			return left.Score > right.Score
		}
		if statusRank(left.Status) != statusRank(right.Status) {
			return statusRank(left.Status) > statusRank(right.Status)
		}
		if left.UpdatedAt != right.UpdatedAt {
			return left.UpdatedAt > right.UpdatedAt
		}
		return left.ID < right.ID
	})
}

func rankCandidates(candidates []Candidate) {
	for i := range candidates {
		candidates[i].Rank = rankCandidate(candidates[i])
	}
	sortCandidates(candidates)
}

func Recommend(records []model.Record, request Request) []Candidate {
	candidates := BuildCandidates(records, request)
	rankCandidates(candidates)
	request = NormalizeRequest(request)
	if len(candidates) > request.Limit {
		candidates = candidates[:request.Limit]
	}
	return candidates
}
