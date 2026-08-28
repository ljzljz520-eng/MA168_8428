package clock

import "testing"

func TestFixedClock(t *testing.T) {
	c := NewFixed(10)
	if c.Now() != 10 || c.Advance(5) != 15 {
		t.Fatal("clock failed")
	}
	c.Set(2)
	if c.Now() != 2 {
		t.Fatal("set failed")
	}
	if len(Sequence(1, 2, 3)) != 3 {
		t.Fatal("sequence failed")
	}
}
