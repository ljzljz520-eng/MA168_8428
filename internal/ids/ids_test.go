package ids

import "testing"

func TestGeneratorAndNormalize(t *testing.T) {
	g := New("audit")
	if g.Peek() != "audit-0001" || g.Next() != "audit-0001" || g.Next() != "audit-0002" {
		t.Fatal("generator failed")
	}
	if Normalize(" Hello/Book ") != "hellobook" {
		t.Fatal("normalize failed")
	}
	if Join("Store", " 1") != "store-1" {
		t.Fatal("join failed")
	}
}
