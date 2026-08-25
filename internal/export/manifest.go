package export

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type Manifest struct {
	Format string `json:"format"`
	Count  int    `json:"count"`
	Digest string `json:"digest"`
}

func BuildManifest(records []byte, count int) Manifest {
	digest := sha256.Sum256(records)
	return Manifest{Format: "bookstore-recommendation-v1", Count: count, Digest: hex.EncodeToString(digest[:])}
}

func EncodeManifest(manifest Manifest) []byte { data, _ := json.Marshal(manifest); return data }

func VerifyManifest(records []byte, manifest Manifest) bool {
	return BuildManifest(records, manifest.Count).Digest == manifest.Digest
}

func MergeManifests(left, right Manifest) Manifest {
	if left.Format == "" {
		return right
	}
	if right.Format == "" {
		return left
	}
	return Manifest{Format: left.Format, Count: left.Count + right.Count, Digest: left.Digest + right.Digest}
}

func IsCompatible(manifest Manifest) bool {
	return manifest.Format == "bookstore-recommendation-v1" && manifest.Count >= 0 && len(manifest.Digest) == 64
}
