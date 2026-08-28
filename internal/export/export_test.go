package export

import (
	"bookstore/recommendation/internal/model"
	"bytes"
	"testing"
)

func TestExportBundleAndManifest(t *testing.T) {
	var output bytes.Buffer
	records := []model.Record{{ID: "r", StoreID: "s", Title: "Book", Score: 80, Status: model.StatusApproved}}
	if err := WriteCSV(&output, records); err != nil || !bytes.Contains(output.Bytes(), []byte("Book")) {
		t.Fatal("csv export failed")
	}
	manifest := BuildManifest(output.Bytes(), 1)
	if !VerifyManifest(output.Bytes(), manifest) {
		t.Fatal("manifest verification failed")
	}
	var jsonOutput bytes.Buffer
	if err := WriteJSON(&jsonOutput, NewBundle(records, nil, nil)); err != nil || jsonOutput.Len() == 0 {
		t.Fatal("json export failed")
	}
}
