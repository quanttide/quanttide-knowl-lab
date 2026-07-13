package main

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func TestDetectBylaw(t *testing.T) {
	dir := "../../_examples/bylaw"
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := runDetect(fs, []string{"-input", dir}); err != nil {
		t.Fatal(err)
	}
}

func TestDetectInvalidDir(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := runDetect(fs, []string{"-input", "/nonexistent"}); err != nil {
		t.Fatal(err)
	}
}

func TestDetectNoInput(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := runDetect(fs, []string{}); err == nil {
		t.Error("expected error for missing -input")
	}
}

func TestValidateBylaw(t *testing.T) {
	dir := "../../_examples/bylaw"
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := runValidate(fs, []string{"-input", dir}); err != nil {
		t.Fatal(err)
	}
}

func TestValidateNoInput(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := runValidate(fs, []string{}); err == nil {
		t.Error("expected error for missing -input")
	}
}

func TestTranscodeBylaw(t *testing.T) {
	td := t.TempDir()
	input := "../../_examples/bylaw"
	output := filepath.Join(td, "archive")

	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := runTranscode(fs, []string{"-input", input, "-output", output}); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(output, "archive.jsonl")); err != nil {
		t.Errorf("expected archive.jsonl: %v", err)
	}
}

func TestTranscodeNoInput(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := runTranscode(fs, []string{}); err == nil {
		t.Error("expected error")
	}
}

func TestTranscodeNoOutput(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := runTranscode(fs, []string{"-input", "."}); err == nil {
		t.Error("expected error")
	}
}

func TestValidateBylawFull(t *testing.T) {
	dir := "../../_examples/bylaw"
	fs := flag.NewFlagSet("test", flag.PanicOnError)
	if err := runValidate(fs, []string{"-input", dir}); err != nil {
		t.Fatal(err)
	}
}
