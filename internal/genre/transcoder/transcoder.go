package transcoder

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/genre"
)

type ArchiveEntry struct {
	Genre     string `json:"genre"`
	Source    string `json:"source"`
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
}

type Transcoder struct{}

func New() *Transcoder {
	return &Transcoder{}
}

func (t *Transcoder) ManualToArchive(input, output string, m *genre.Metadata) (string, error) {
	if err := os.MkdirAll(output, 0755); err != nil {
		return "", fmt.Errorf("creating output dir: %w", err)
	}

	knowlDir := filepath.Join(input, "knowl")
	entries, err := os.ReadDir(knowlDir)
	if err != nil {
		return "", fmt.Errorf("reading knowl dir: %w", err)
	}

	var archive []ArchiveEntry
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(knowlDir, e.Name()))
		if err != nil {
			return "", fmt.Errorf("reading %s: %w", e.Name(), err)
		}
		lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			var entry struct {
				Subject   string  `json:"Subject"`
				Predicate string  `json:"Predicate"`
				Object    string  `json:"Object"`
				Confidence float64 `json:"Confidence"`
				Source    string  `json:"Source"`
			}
			if err := json.Unmarshal([]byte(line), &entry); err != nil {
				continue
			}
			archive = append(archive, ArchiveEntry{
				Genre:     string(genre.Archive),
				Source:    entry.Source,
				Subject:   entry.Subject,
				Predicate: entry.Predicate,
				Object:    entry.Object,
			})
		}
	}

	outPath := filepath.Join(output, "archive.jsonl")
	out, err := os.Create(outPath)
	if err != nil {
		return "", fmt.Errorf("creating output file: %w", err)
	}
	defer out.Close()

	enc := json.NewEncoder(out)
	for _, a := range archive {
		if err := enc.Encode(a); err != nil {
			return "", fmt.Errorf("encoding entry: %w", err)
		}
	}

	return outPath, nil
}
