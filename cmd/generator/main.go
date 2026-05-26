package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/generator"
	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/llm"
	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/rule"
	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/storage"
)

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if err := run(fs, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fs *flag.FlagSet, args []string) error {
	target := fs.String("target", "code", "generator target (code / api / faq)")
	query := fs.String("query", "", "query predicate (e.g. contains)")
	output := fs.String("output", "./generated", "output directory or file")
	store := fs.String("store", "./knowl/store.jsonl", "path to store.jsonl")
	apiURL := fs.String("api-url", "", "LLM API URL (optional, enables LLM-enhanced FAQ)")
	model := fs.String("model", "deepseek-chat", "LLM model name")
	outFormat := fs.String("format", "markdown", "FAQ output format (markdown / sharegpt / alpaca / reasoning)")
	eval := fs.Bool("eval", false, "evaluate FAQ quality against source triples")
	evalTriples := fs.String("eval-triples", "./knowl/triples.jsonl", "source triples for evaluation")
	fs.Parse(args)

	if *eval {
		if *apiURL == "" {
			return fmt.Errorf("-eval 需要 -api-url 和 DEEPSEEK_API_KEY")
		}
		client := llm.New(*apiURL, os.Getenv("DEEPSEEK_API_KEY"), *model)
		return runEval(client, *output, *evalTriples)
	}

	gen, err := generator.GetGenerator(*target)
	if err != nil {
		return err
	}

	if *target == "faq" {
		if faqGen, ok := gen.(*generator.FAQGenerator); ok {
			faqGen.SetOutFormat(*outFormat)
			if *apiURL != "" {
				c := llm.New(*apiURL, os.Getenv("DEEPSEEK_API_KEY"), *model)
				faqGen.SetLLMClient(c)
			}
		}
	}

	facts, err := loadFacts(*store, *query)
	if err != nil {
		return err
	}
	if len(facts) < generator.MinFacts {
		return fmt.Errorf("仅 %d 条事实，需要至少 %d 条才能生成", len(facts), generator.MinFacts)
	}

	result, err := gen.Generate(facts, *query)
	if err != nil {
		return err
	}

	if err := writeOutput(result, *output); err != nil {
		fmt.Fprintf(os.Stderr, "warn: cannot write to output path: %v\n", err)
		fmt.Println("--- stdout ---")
		fmt.Println(result.Content)
		return nil
	}

	if result.Format == "go" {
		if err := gofmt(result.Path); err != nil {
			fmt.Fprintf(os.Stderr, "warn: go fmt failed (output preserved): %v\n", err)
		}
	}

	fmt.Printf("✓ 已生成 %s 到 %s（基于 %d 条事实）\n", result.Type, result.Path, len(facts))
	return nil
}

func loadFacts(storePath, query string) ([]generator.Fact, error) {
	predicate := query
	if predicate == "*" {
		predicate = ""
	} else if strings.Contains(query, "(") {
		atom, err := rule.ParseAtom(query)
		if err == nil {
			predicate = atom.Predicate
		}
	}

	f, err := os.Open(storePath)
	if err != nil {
		return nil, fmt.Errorf("打开 store 文件: %w", err)
	}
	defer f.Close()

	var facts []generator.Fact
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var rec storage.TripleRecord
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			continue
		}
		if predicate != "" && rec.Predicate != predicate {
			continue
		}
		facts = append(facts, generator.Fact{
			Subject:   rec.Subject,
			Predicate: rec.Predicate,
			Object:    rec.Object,
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return facts, nil
}

func writeOutput(result *generator.Output, outputPath string) error {
	info, err := os.Stat(outputPath)
	isDir := err == nil && info.IsDir()

	var filePath string
	switch {
	case isDir:
		switch result.Format {
		case "go":
			filePath = filepath.Join(outputPath, "types.go")
		case "markdown":
			filePath = filepath.Join(outputPath, result.Type+".md")
		default:
			filePath = filepath.Join(outputPath, result.Type+".txt")
		}
	case strings.HasSuffix(outputPath, "/") || outputPath == "":
		return fmt.Errorf("invalid output path: %s", outputPath)
	default:
		filePath = outputPath
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	if err := os.WriteFile(filePath, []byte(result.Content), 0644); err != nil {
		return err
	}

	result.Path = filePath
	return nil
}

func gofmt(path string) error {
	cmd := exec.Command("go", "fmt", path)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go fmt: %w\n%s", err, string(output))
	}
	return nil
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "用法: go run ./cmd/generator -target <code|api|faq> -query <predicate> [-output <path>] [-store <store.jsonl>]\n\n")
		fmt.Fprintf(os.Stderr, "可用生成器:\n")
		for name := range generator.AvailableGenerators {
			fmt.Fprintf(os.Stderr, "  %s\n", name)
		}
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  go run ./cmd/generator -target faq -query 是 -output ./wisdom/faq.md\n")
		fmt.Fprintf(os.Stderr, "  go run ./cmd/generator -target code -query 包含 -output ./wisdom/types.go\n")
	}
}
