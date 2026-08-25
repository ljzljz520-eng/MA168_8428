package workflow

import (
	"bookstore/recommendation/internal/store"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkflowImportReport(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "import.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	service := New(st)
	input := "id,store,title,author,genre,score,summary\ni1,s,One,A,Fiction,80,first\ni2,s,Two,B,History,70,second\nbad,line\n"
	result := service.ImportCSV(strings.NewReader(input), "importer", 10)
	if len(result.Created) != 2 || len(result.Failed) != 1 {
		t.Fatalf("import result mismatch: %#v", result)
	}
	report, err := service.BuildReport("")
	if err != nil || report.Summary.Total != 2 || report.Audits != 2 {
		t.Fatalf("report mismatch: %v %#v", err, report)
	}
}
