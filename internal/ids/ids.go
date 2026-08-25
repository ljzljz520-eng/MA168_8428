package ids

import (
	"fmt"
	"strings"
)

type Generator struct {
	prefix string
	next   int
}

func New(prefix string) *Generator { return &Generator{prefix: strings.TrimSpace(prefix)} }

func (g *Generator) Next() string { g.next++; return fmt.Sprintf("%s-%04d", g.prefix, g.next) }

func (g *Generator) Peek() string { return fmt.Sprintf("%s-%04d", g.prefix, g.next+1) }

func Normalize(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	result := ""
	for _, r := range value {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' {
			result += string(r)
		}
	}
	return result
}

func Join(parts ...string) string {
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if value := Normalize(part); value != "" {
			result = append(result, value)
		}
	}
	return strings.Join(result, "-")
}
