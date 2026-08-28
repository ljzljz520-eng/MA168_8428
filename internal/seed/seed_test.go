package seed

import (
	"bookstore/recommendation/internal/fixture"
	"bookstore/recommendation/internal/store"
	"bookstore/recommendation/internal/workflow"
	"path/filepath"
	"testing"
)

func TestLoadStandard(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "seed.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	result := New(workflow.New(st)).LoadStandard("seed", 1)
	if err = Validate(result, len(fixture.Standard())); err != nil {
		t.Fatal(err)
	}
}
