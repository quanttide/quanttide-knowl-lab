package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/genre"
)

type Level string

const (
	LevelPass    Level = "PASS"
	LevelWarning Level = "WARN"
	LevelError   Level = "ERROR"
)

type Item struct {
	Level   Level  `json:"level"`
	Message string `json:"message"`
}

type Report struct {
	Genre    genre.Genre `json:"genre"`
	Passed   int         `json:"passed"`
	Warnings int         `json:"warnings"`
	Errors   int         `json:"errors"`
	Items    []Item      `json:"items"`
}

type RuleSet struct {
	Rules []Rule
}

type Rule interface {
	Name() string
	Check(path string) []Item
}

func ForGenre(g genre.Genre) *RuleSet {
	switch g {
	case genre.Manual:
		return &RuleSet{Rules: []Rule{
			&StructureRule{},
			&FieldRule{},
		}}
	default:
		return &RuleSet{}
	}
}

func (rs *RuleSet) Validate(path string) *Report {
	r := &Report{Items: []Item{}}
	for _, rule := range rs.Rules {
		for _, item := range rule.Check(path) {
			r.Items = append(r.Items, item)
			switch item.Level {
			case LevelPass:
				r.Passed++
			case LevelWarning:
				r.Warnings++
			case LevelError:
				r.Errors++
			}
		}
	}
	return r
}

type StructureRule struct{}

func (r *StructureRule) Name() string { return "structure" }

func (r *StructureRule) Check(path string) []Item {
	required := []string{"data", "info", "knowl", "wisdom"}
	var items []Item
	for _, d := range required {
		info, err := os.Stat(filepath.Join(path, d))
		if err != nil || !info.IsDir() {
			items = append(items, Item{Level: LevelError, Message: "missing required directory: " + d})
		} else {
			items = append(items, Item{Level: LevelPass, Message: "directory exists: " + d})
		}
	}
	return items
}

type FieldRule struct{}

func (r *FieldRule) Name() string { return "fields" }

func (r *FieldRule) Check(path string) []Item {
	var items []Item

	wisdomFile := filepath.Join(path, "wisdom", "faq.md")
	raw, err := os.ReadFile(wisdomFile)
	if err != nil {
		return []Item{{Level: LevelError, Message: "cannot read wisdom/faq.md: " + err.Error()}}
	}
	content := string(raw)
	if !strings.Contains(content, "<question>") || !strings.Contains(content, "<answer>") {
		items = append(items, Item{Level: LevelWarning, Message: "wisdom/faq.md may lack question/answer tags"})
	} else {
		items = append(items, Item{Level: LevelPass, Message: "wisdom/faq.md contains question/answer pairs"})
	}

	storeFile := filepath.Join(path, "knowl", "store.jsonl")
	storeRaw, err := os.ReadFile(storeFile)
	if err != nil {
		items = append(items, Item{Level: LevelWarning, Message: "cannot read knowl/store.jsonl"})
	} else {
		lines := strings.Split(strings.TrimSpace(string(storeRaw)), "\n")
		if len(lines) > 0 && lines[0] != "" {
			items = append(items, Item{Level: LevelPass, Message: "knowl/store.jsonl has " + fmt.Sprintf("%d", len(lines)) + " entries"})
		}
	}

	if len(items) == 0 {
		items = append(items, Item{Level: LevelPass, Message: "all required fields present"})
	}
	return items
}
