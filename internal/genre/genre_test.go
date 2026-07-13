package genre

import (
	"os"
	"path/filepath"
	"testing"
)

func TestClassifyBylaw(t *testing.T) {
	result := Classify("../../_examples/bylaw")
	if result.Genre != Manual {
		t.Errorf("expected Manual, got %s", result.Genre)
	}
	if result.Confidence < 0.9 {
		t.Errorf("expected confidence >= 0.9, got %f", result.Confidence)
	}
}

func TestClassifyUnknown(t *testing.T) {
	dir := t.TempDir()
	result := Classify(dir)
	if result.Genre != Unknown {
		t.Errorf("expected Unknown, got %s", result.Genre)
	}
}

func TestClassifyWithGenreJSON(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "data"), 0755)
	os.MkdirAll(filepath.Join(dir, "info"), 0755)
	os.MkdirAll(filepath.Join(dir, "knowl"), 0755)
	os.MkdirAll(filepath.Join(dir, "wisdom"), 0755)
	os.WriteFile(filepath.Join(dir, "genre.json"), []byte(`{"genre":"Manual","contract":{"form":"手册","granularity":"章节","lifecycle_stages":["编写","评审","发布","修订"]}}`), 0644)

	result := Classify(dir)
	if result.Genre != Manual {
		t.Errorf("expected Manual, got %s", result.Genre)
	}
	if result.Confidence != 1.0 {
		t.Errorf("expected confidence 1.0, got %f", result.Confidence)
	}
}

func TestLoadMetadata(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "genre.json"), []byte(`{"genre":"Manual","contract":{"form":"手册","granularity":"章节","lifecycle_stages":["编写","评审"]}}`), 0644)

	m, err := LoadMetadata(dir)
	if err != nil {
		t.Fatal(err)
	}
	if m.Genre != Manual {
		t.Errorf("expected Manual, got %s", m.Genre)
	}
	if m.Contract.Form != "手册" {
		t.Errorf("expected 手册, got %s", m.Contract.Form)
	}
}

func TestLoadMetadataMissing(t *testing.T) {
	dir := t.TempDir()
	_, err := LoadMetadata(dir)
	if err == nil {
		t.Error("expected error for missing genre.json")
	}
}
