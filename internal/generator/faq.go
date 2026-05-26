package generator

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/llm"
)

type FAQGenerator struct {
	client   *llm.Client
	OutFormat string // markdown, sharegpt, alpaca, reasoning
}

func (g *FAQGenerator) Name() string { return "faq" }

func (g *FAQGenerator) SetLLMClient(c *llm.Client) {
	g.client = c
}

func (g *FAQGenerator) SetOutFormat(f string) {
	g.OutFormat = f
}

func (g *FAQGenerator) Generate(facts []Fact, query string) (*Output, error) {
	if len(facts) < MinFacts {
		return nil, fmt.Errorf("数据不足：仅 %d 条事实，需要至少 %d 条才能生成 FAQ", len(facts), MinFacts)
	}

	if g.client != nil {
		return g.generateWithLLM(facts)
	}
	return g.generateWithTemplate(facts), nil
}

func (g *FAQGenerator) generateWithLLM(facts []Fact) (*Output, error) {
	pairs, err := g.generateQAPairs(facts)
	if err != nil {
		return nil, err
	}

	content := formatFAQ(pairs, g.OutFormat)
	ext := "md"
	if g.OutFormat == "sharegpt" || g.OutFormat == "alpaca" {
		ext = "jsonl"
	}

	return &Output{
		Type:    "faq",
		Content: content,
		Format:  ext,
	}, nil
}

func (g *FAQGenerator) generateQAPairs(facts []Fact) ([]qaPair, error) {
	var triplesJSON strings.Builder
	for _, f := range facts {
		triplesJSON.WriteString(fmt.Sprintf("  - %s → %s → %s\n", f.Subject, f.Predicate, f.Object))
	}

	prompt := fmt.Sprintf(`你是一位知识库编辑。基于以下三元组列表，为每个三元组生成一个问答对。

每条三元组表示 (主语, 谓语, 宾语)。

要求：
- 语言自然流畅，符合中文表达习惯
- 如果谓语是"是"、"属于"、"即"，问"主语是什么？"
- 如果谓语是其他词，问"主语的谓语是什么？"
- 输出格式：每行一个 JSON 对象：{"question": "...", "answer": "..."}

三元组列表：
%s`, triplesJSON.String())

	resp, err := g.client.Chat([]llm.Message{
		{Role: "system", Content: "你是一个专业的知识库编辑，擅长将结构化数据转化为问答对。输出 JSON，每行一个。"},
		{Role: "user", Content: prompt},
	})
	if err != nil {
		return nil, fmt.Errorf("LLM 生成 FAQ 失败: %w", err)
	}

	return parseQAPairs(resp.Choices[0].Message.Content)
}

type qaPair struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func parseQAPairs(raw string) ([]qaPair, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	var pairs []qaPair
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var p qaPair
		if err := json.Unmarshal([]byte(line), &p); err != nil {
			continue
		}
		if p.Question != "" && p.Answer != "" {
			pairs = append(pairs, p)
		}
	}
	if len(pairs) == 0 {
		return nil, fmt.Errorf("未能从 LLM 响应中解析出有效问答对")
	}
	return pairs, nil
}

func formatFAQ(pairs []qaPair, format string) string {
	switch format {
	case "sharegpt":
		return formatShareGPT(pairs)
	case "alpaca":
		return formatAlpaca(pairs)
	case "reasoning":
		return formatReasoning(pairs)
	default:
		return formatMarkdown(pairs)
	}
}

func formatMarkdown(pairs []qaPair) string {
	var buf strings.Builder
	buf.WriteString("# FAQ - 常见问题\n\n")
	for i, p := range pairs {
		buf.WriteString(fmt.Sprintf("## Q%d: %s\n\n%s\n\n", i+1, p.Question, p.Answer))
	}
	return buf.String()
}

func formatShareGPT(pairs []qaPair) string {
	var buf strings.Builder
	for _, p := range pairs {
		data, _ := json.Marshal(map[string][]map[string]string{
			"conversations": {
				{"from": "human", "value": p.Question},
				{"from": "gpt", "value": p.Answer},
			},
		})
		buf.WriteString(string(data))
		buf.WriteString("\n")
	}
	return buf.String()
}

func formatAlpaca(pairs []qaPair) string {
	var buf strings.Builder
	for _, p := range pairs {
		data, _ := json.Marshal(map[string]string{
			"instruction": "基于以下知识回答问题",
			"input":       p.Question,
			"output":      p.Answer,
		})
		buf.WriteString(string(data))
		buf.WriteString("\n")
	}
	return buf.String()
}

func formatReasoning(pairs []qaPair) string {
	var buf strings.Builder
	for _, p := range pairs {
		buf.WriteString(fmt.Sprintf("<question>%s</question>\n<answer>%s</answer>\n\n", p.Question, p.Answer))
	}
	return buf.String()
}

func (g *FAQGenerator) generateWithTemplate(facts []Fact) *Output {
	var faq strings.Builder
	faq.WriteString("# FAQ - 常见问题\n\n")
	faq.WriteString(fmt.Sprintf("> 基于 %d 条知识事实自动生成\n\n", len(facts)))

	seen := make(map[string]bool)
	for i, f := range facts {
		q := templateQuestion(f)
		if seen[q] {
			continue
		}
		seen[q] = true
		a := templateAnswer(f)
		faq.WriteString(fmt.Sprintf("## Q%d: %s\n\n%s\n\n", i+1, q, a))
	}

	return &Output{
		Type:    "faq",
		Content: faq.String(),
		Format:  "markdown",
	}
}

var bePredicates = map[string]bool{"是": true, "属于": true, "即": true}

func templateQuestion(f Fact) string {
	if bePredicates[f.Predicate] {
		return fmt.Sprintf("%s 是什么？", f.Subject)
	}
	return fmt.Sprintf("%s 的 %s 是什么？", f.Subject, f.Predicate)
}

func templateAnswer(f Fact) string {
	if bePredicates[f.Predicate] {
		return fmt.Sprintf("%s 是 %s。", f.Subject, f.Object)
	}
	return fmt.Sprintf("%s 的 %s 是 %s。", f.Subject, f.Predicate, f.Object)
}


