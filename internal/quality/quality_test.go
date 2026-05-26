package quality

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeEmptyDir(t *testing.T) {
	r, docs, err := AnalyzeDir(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.DocCount != 0 {
		t.Errorf("expected 0 docs, got %d", r.DocCount)
	}
	if len(docs) != 0 {
		t.Errorf("expected 0 docs, got %d", len(docs))
	}
}

func TestAnalyzeSingleDoc(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.md"), []byte("# Title\n\nSome content."), 0644)

	r, docs, err := AnalyzeDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.DocCount != 1 {
		t.Errorf("expected 1 doc, got %d", r.DocCount)
	}
	if len(docs) != 1 {
		t.Errorf("expected 1 doc, got %d", len(docs))
	}
	if docs[0].HeadingCount != 1 {
		t.Errorf("expected 1 heading, got %d", docs[0].HeadingCount)
	}
}

func TestAnalyzeWithFrontmatter(t *testing.T) {
	dir := t.TempDir()
	content := "---\ntitle: Test\n---\n\n# Section\n\nBody text."
	os.WriteFile(filepath.Join(dir, "fm.md"), []byte(content), 0644)

	_, docs, err := AnalyzeDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !docs[0].HasFrontmatter {
		t.Error("expected frontmatter detection")
	}
}

func TestAnalyzeListsAndTables(t *testing.T) {
	dir := t.TempDir()
	content := "# Report\n\n- item 1\n- item 2\n- item 3\n\n| col1 | col2 |\n| a | b |\n"
	os.WriteFile(filepath.Join(dir, "rpt.md"), []byte(content), 0644)

	_, docs, err := AnalyzeDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if docs[0].ListCount < 2 {
		t.Errorf("expected at least 2 list items, got %d", docs[0].ListCount)
	}
	if docs[0].TableCount < 1 {
		t.Errorf("expected at least 1 table row, got %d", docs[0].TableCount)
	}
}

func TestRecommendationEmpty(t *testing.T) {
	r := &Report{DocCount: 0}
	v, _ := Recommendation(r)
	if v != "skip" {
		t.Errorf("expected skip, got %s", v)
	}
}

func TestRecommendationProceed(t *testing.T) {
	r := &Report{
		DocCount:       5,
		StructureScore: 60,
		AvgLength:      500,
		AvgHeadingCount: 3,
		HasFrontmatterPct: 80,
	}
	v, _ := Recommendation(r)
	if v != "proceed" {
		t.Errorf("expected proceed, got %s", v)
	}
}

func TestRecommendationReview(t *testing.T) {
	r := &Report{
		DocCount:       1,
		StructureScore: 20,
		AvgLength:      100,
		AvgHeadingCount: 0,
	}
	v, _ := Recommendation(r)
	if v != "review" {
		t.Errorf("expected review, got %s", v)
	}
}

func TestAnalyzeMultipleDocs(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "a.md"), []byte("# A\n\nContent."), 0644)
	os.WriteFile(filepath.Join(dir, "b.txt"), []byte("Plain text"), 0644)
	os.WriteFile(filepath.Join(dir, "c.json"), []byte("{}"), 0644)

	r, docs, err := AnalyzeDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.DocCount != 2 {
		t.Errorf("expected 2 docs (.md + .txt), got %d", r.DocCount)
	}
	if len(docs) != 2 {
		t.Errorf("expected 2 docs, got %d", len(docs))
	}
}

func TestAvgParaLength(t *testing.T) {
	dir := t.TempDir()
	content := "# H\n\nLong paragraph with more words for testing average calculation.\n\nAnother short one."
	os.WriteFile(filepath.Join(dir, "p.md"), []byte(content), 0644)

	_, docs, err := AnalyzeDir(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if docs[0].AvgParaLen <= 0 {
		t.Errorf("expected positive avg paragraph length, got %f", docs[0].AvgParaLen)
	}
}
