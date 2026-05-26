package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/llm"
)

type evalTriple struct {
	Subject   string `json:"subject"`
	Predicate string `json:"predicate"`
	Object    string `json:"object"`
	Sentence  string `json:"sentence"`
}

type evalResponse struct {
	Faithfulness int `json:"faithfulness"`
	Usefulness   int `json:"usefulness"`
	Completeness int `json:"completeness"`
}

type evalResult struct {
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	Sentence     string `json:"sentence"`
	Faithfulness int    `json:"faithfulness"`
	Usefulness   int    `json:"usefulness"`
	Completeness int    `json:"completeness"`
}

func runEval(client *llm.Client, faqPath, triplesPath string) error {
	faqs, err := parseFAQs(faqPath)
	if err != nil {
		return fmt.Errorf("parse faq: %w", err)
	}

	triples, err := parseTriplesWithContext(triplesPath)
	if err != nil {
		return fmt.Errorf("parse triples: %w", err)
	}

	var results []evalResult
	totalFaith, totalUse, totalComp := 0, 0, 0

	for idx, faq := range faqs {
		tc := matchSentence(faq.answer, triples)
		if tc == "" {
			tc = "（未找到匹配源句）"
		}

		prompt := fmt.Sprintf(`你是一个 FAQ 质量评估员。评估以下问答对：

问：%s
答：%s

来源原文：%s

请从 1-5 分评估三个维度：
- faithfulness（忠实度）：答案是否准确反映原文？1=编造, 5=准确
- usefulness（有用度）：这个问答本身是否有信息价值？1=废话, 5=干货
- completeness（完整度）：答案是否遗漏重要条件？1=缺失关键, 5=完整

输出 JSON：{"faithfulness": N, "usefulness": N, "completeness": N}`, faq.question, faq.answer, tc)

		resp, err := client.Chat([]llm.Message{
			{Role: "system", Content: "你是一个严格的质量评估员。输出 JSON。"},
			{Role: "user", Content: prompt},
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: 条目 %d 评估失败: %v\n", idx+1, err)
			continue
		}

		var score evalResponse
		if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &score); err != nil {
			fmt.Fprintf(os.Stderr, "warn: 条目 %d 解析评分失败: %v\n", idx+1, err)
			continue
		}

		results = append(results, evalResult{
			Question:     faq.question,
			Answer:       faq.answer,
			Sentence:     tc,
			Faithfulness: score.Faithfulness,
			Usefulness:   score.Usefulness,
			Completeness: score.Completeness,
		})
		totalFaith += score.Faithfulness
		totalUse += score.Usefulness
		totalComp += score.Completeness
	}

	n := len(results)
	if n == 0 {
		return fmt.Errorf("没有成功评估任何条目")
	}

	fmt.Println("\n=== FAQ 质量报告 ===")
	fmt.Printf("样本数: %d\n", n)
	fmt.Printf("faithfulness: %.1f/5\n", float64(totalFaith)/float64(n))
	fmt.Printf("usefulness:   %.1f/5\n", float64(totalUse)/float64(n))
	fmt.Printf("completeness: %.1f/5\n", float64(totalComp)/float64(n))
	fmt.Println()

	for i, r := range results {
		fmt.Printf("[%d] 忠实=%d 有用=%d 完整=%d | %s\n",
			i+1, r.Faithfulness, r.Usefulness, r.Completeness, r.Question)
	}

	return nil
}

type faqEntry struct {
	question string
	answer   string
}

func parseFAQs(path string) ([]faqEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`<question>(.*?)</question>\s*<answer>(.*?)</answer>`)
	matches := re.FindAllStringSubmatch(string(data), -1)

	var entries []faqEntry
	for _, m := range matches {
		if len(m) >= 3 {
			entries = append(entries, faqEntry{
				question: strings.TrimSpace(m[1]),
				answer:   strings.TrimSpace(m[2]),
			})
		}
	}
	return entries, nil
}

func parseTriplesWithContext(path string) ([]evalTriple, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var triples []evalTriple
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var t struct {
			Subject   string `json:"subject"`
			Predicate string `json:"predicate"`
			Object    string `json:"object"`
			Context   map[string]string `json:"context,omitempty"`
		}
		if err := json.Unmarshal([]byte(line), &t); err != nil {
			continue
		}
		sentence := ""
		if t.Context != nil {
			sentence = t.Context["sentence"]
		}
		triples = append(triples, evalTriple{
			Subject:   t.Subject,
			Predicate: t.Predicate,
			Object:    t.Object,
			Sentence:  sentence,
		})
	}
	return triples, scanner.Err()
}

func matchSentence(answer string, triples []evalTriple) string {
	for _, t := range triples {
		if strings.Contains(answer, t.Object) || strings.Contains(t.Object, answer) {
			if t.Sentence != "" {
				return t.Sentence
			}
			return fmt.Sprintf("%s → %s → %s", t.Subject, t.Predicate, t.Object)
		}
	}
	return ""
}
