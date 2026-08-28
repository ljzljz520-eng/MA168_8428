package recommendation

import "bookstore/recommendation/internal/model"

type Collection struct {
	Items       []Candidate `json:"items"`
	GeneratedAt int64       `json:"generated_at"`
}

func NewCollection(candidates []Candidate, generatedAt int64) Collection {
	items := make([]Candidate, len(candidates))
	copy(items, candidates)
	return Collection{Items: items, GeneratedAt: generatedAt}
}

func (c Collection) IDs() []string {
	ids := make([]string, 0, len(c.Items))
	for _, item := range c.Items {
		ids = append(ids, item.Record.ID)
	}
	return ids
}

func (c Collection) Visible() []model.Record {
	result := make([]model.Record, 0, len(c.Items))
	for _, item := range c.Items {
		if item.Record.IsVisible() {
			result = append(result, item.Record)
		}
	}
	return result
}
