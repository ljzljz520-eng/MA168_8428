package fixture

import "testing"

func TestStandardFixture(t *testing.T) {
	specs := Standard()
	if len(specs) < 4 || len(Inputs(specs)) != len(specs) {
		t.Fatal("fixture incomplete")
	}
	if _, ok := Lookup(specs, "fiction-01"); !ok {
		t.Fatal("lookup failed")
	}
	if len(CSV(specs)) == 0 {
		t.Fatal("csv empty")
	}
}
