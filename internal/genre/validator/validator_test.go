package validator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/genre"
)

func TestManualStructureValid(t *testing.T) {
	dir := t.TempDir()
	for _, d := range []string{"data", "info", "knowl", "wisdom"} {
		os.MkdirAll(filepath.Join(dir, d), 0755)
	}
	os.WriteFile(filepath.Join(dir, "wisdom", "faq.md"), []byte("<question>q</question>\n<answer>a</answer>"), 0644)

	rs := ForGenre(genre.Manual)
	report := rs.Validate(dir)
	if report.Errors > 0 {
		t.Errorf("expected 0 errors, got %d", report.Errors)
	}
}

func TestManualStructureMissing(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, "data"), 0755)

	rs := ForGenre(genre.Manual)
	report := rs.Validate(dir)
	if report.Errors == 0 {
		t.Error("expected errors for missing directories")
	}
}

func TestManualFieldsValid(t *testing.T) {
	dir := t.TempDir()
	for _, d := range []string{"data", "info", "knowl", "wisdom"} {
		os.MkdirAll(filepath.Join(dir, d), 0755)
	}
	os.WriteFile(filepath.Join(dir, "wisdom", "faq.md"), []byte("<question>q</question>\n<answer>a</answer>"), 0644)
	os.WriteFile(filepath.Join(dir, "knowl", "store.jsonl"), []byte("{\"key\":\"val\"}\n"), 0644)

	rs := ForGenre(genre.Manual)
	report := rs.Validate(dir)
	if report.Errors > 0 {
		t.Errorf("expected 0 errors, got %d", report.Errors)
	}
}

func TestManualFieldsMissingFAQ(t *testing.T) {
	dir := t.TempDir()
	for _, d := range []string{"data", "info", "knowl", "wisdom"} {
		os.MkdirAll(filepath.Join(dir, d), 0755)
	}
	os.WriteFile(filepath.Join(dir, "wisdom", "faq.md"), []byte("no tags here"), 0644)
	os.WriteFile(filepath.Join(dir, "knowl", "store.jsonl"), []byte("{}"), 0644)

	rs := ForGenre(genre.Manual)
	report := rs.Validate(dir)
	if report.Warnings == 0 {
		t.Error("expected warnings for missing question/answer tags")
	}
}

func TestManualFieldsMissingStore(t *testing.T) {
	dir := t.TempDir()
	for _, d := range []string{"data", "info", "knowl", "wisdom"} {
		os.MkdirAll(filepath.Join(dir, d), 0755)
	}
	os.WriteFile(filepath.Join(dir, "wisdom", "faq.md"), []byte("<question>q</question>\n<answer>a</answer>"), 0644)

	rs := ForGenre(genre.Manual)
	report := rs.Validate(dir)
	if report.Warnings == 0 {
		t.Error("expected warnings for missing store.jsonl fields")
	}
}

func TestReportCounts(t *testing.T) {
	dir := t.TempDir()
	for _, d := range []string{"data", "info", "knowl", "wisdom"} {
		os.MkdirAll(filepath.Join(dir, d), 0755)
	}
	os.WriteFile(filepath.Join(dir, "wisdom", "faq.md"), []byte("<question>q</question>\n<answer>a</answer>"), 0644)
	os.WriteFile(filepath.Join(dir, "knowl", "store.jsonl"), []byte("{}"), 0644)

	rs := ForGenre(genre.Manual)
	report := rs.Validate(dir)
	if report.Passed == 0 {
		t.Error("expected at least 1 passed check")
	}
}
