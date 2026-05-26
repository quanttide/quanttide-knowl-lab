package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func writeTestStore(t *testing.T, dir string) string {
	store := filepath.Join(dir, "store.jsonl")
	content := `{"Subject":"用户","Predicate":"创建","Object":"订单","Confidence":0.9}
{"Subject":"订单","Predicate":"包含","Object":"商品","Confidence":0.8}
{"Subject":"订单","Predicate":"属于","Object":"用户","Confidence":0.7}
{"Subject":"商品","Predicate":"包含","Object":"库存","Confidence":0.95}
{"Subject":"电子产品","Predicate":"包含","Object":"配件","Confidence":0.85}
{"Subject":"退换货","Predicate":"期限","Object":"30天","Confidence":0.9}`
	if err := os.WriteFile(store, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return store
}

func TestGeneratorCode(t *testing.T) {
	dir := t.TempDir()
	store := writeTestStore(t, dir)
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := run(fs, []string{"-target", "code", "-query", "包含", "-output", dir, "-store", store}); err != nil {
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) == 0 {
		t.Error("expected generated files")
	}
}

func TestGeneratorAPI(t *testing.T) {
	dir := t.TempDir()
	store := writeTestStore(t, dir)
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := run(fs, []string{"-target", "api", "-query", "包含", "-output", dir, "-store", store}); err != nil {
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) == 0 {
		t.Error("expected generated files")
	}
}

func TestGeneratorFAQ(t *testing.T) {
	dir := t.TempDir()
	store := writeTestStore(t, dir)
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := run(fs, []string{"-target", "faq", "-query", "包含", "-output", dir, "-store", store}); err != nil {
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) == 0 {
		t.Error("expected generated files")
	}
}

func TestGeneratorUnknownTarget(t *testing.T) {
	dir := t.TempDir()
	store := writeTestStore(t, dir)
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := run(fs, []string{"-target", "unknown", "-query", "包含", "-output", dir, "-store", store}); err == nil {
		t.Error("expected error for unknown target")
	}
}

func TestGeneratorInsufficientFacts(t *testing.T) {
	dir := t.TempDir()
	store := filepath.Join(dir, "store.jsonl")
	os.WriteFile(store, []byte(`{"Subject":"A","Predicate":"B","Object":"C","Confidence":0.9,"Source":"doc.md"}`), 0644)
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := run(fs, []string{"-target", "faq", "-query", "B", "-output", dir, "-store", store}); err == nil {
		t.Error("expected error for insufficient facts")
	}
}

func TestGeneratorNoStore(t *testing.T) {
	dir := t.TempDir()
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := run(fs, []string{"-target", "code", "-query", "包含", "-output", dir, "-store", "/nonexistent/store.jsonl"}); err == nil {
		t.Error("expected error for missing store")
	}
}

func TestGeneratorOutputToFile(t *testing.T) {
	dir := t.TempDir()
	store := writeTestStore(t, dir)
	out := filepath.Join(dir, "out.go")
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := run(fs, []string{"-target", "code", "-query", "包含", "-output", out, "-store", store}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(out); err != nil {
		t.Errorf("expected output file: %v", err)
	}
}
