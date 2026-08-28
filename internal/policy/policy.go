package policy

import (
	"bookstore/recommendation/internal/model"
	"fmt"
	"strings"
)

type Rule struct {
	ID           string
	Description  string
	MinimumScore int
	RequiredTags []string
}

type Decision struct {
	Allowed bool
	RuleID  string
	Reasons []string
}

type Policy struct {
	rules       map[string]Rule
	defaultRule Rule
}

func New() *Policy {
	return &Policy{rules: map[string]Rule{
		"fiction": {ID: "fiction-quality", Description: "fiction entries need a minimum editorial score", MinimumScore: 60, RequiredTags: []string{}},
		"history": {ID: "history-context", Description: "history entries need a context tag", MinimumScore: 55, RequiredTags: []string{"context"}},
	}, defaultRule: Rule{ID: "general-quality", Description: "general entries need a consistent score", MinimumScore: 50}}
}

func (p *Policy) RuleFor(genre string) Rule {
	if rule, ok := p.rules[strings.ToLower(strings.TrimSpace(genre))]; ok {
		return rule
	}
	return p.defaultRule
}

func (p *Policy) Evaluate(record model.Record) Decision {
	rule := p.RuleFor(record.Genre)
	reasons := []string{}
	if record.Score < rule.MinimumScore {
		reasons = append(reasons, fmt.Sprintf("score %d is below %d", record.Score, rule.MinimumScore))
	}
	for _, required := range rule.RequiredTags {
		if !contains(record.Tags, required) {
			reasons = append(reasons, "missing tag "+required)
		}
	}
	if record.Title == "" || record.Author == "" {
		reasons = append(reasons, "title and author are required")
	}
	return Decision{Allowed: len(reasons) == 0, RuleID: rule.ID, Reasons: reasons}
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(value, target) {
			return true
		}
	}
	return false
}

func (p *Policy) Explain(record model.Record) string {
	decision := p.Evaluate(record)
	if decision.Allowed {
		return decision.RuleID + ": accepted"
	}
	return decision.RuleID + ": " + strings.Join(decision.Reasons, "; ")
}
