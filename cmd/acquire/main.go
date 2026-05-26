package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/llm"
	"github.com/quanttide/quanttide-example-of-knowledge-engineering/internal/quality"
)

func main() {
	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	if err := run(fs, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(fs *flag.FlagSet, args []string) error {
	input := fs.String("input", "./data", "input directory with documents")
	apiURL := fs.String("api-url", "", "LLM API URL (optional, enables semantic assessment)")
	model := fs.String("model", "deepseek-chat", "LLM model name")
	fs.Parse(args)

	r, docs, err := quality.AnalyzeDir(*input)
	if err != nil {
		return fmt.Errorf("analyze: %w", err)
	}

	fmt.Println("=== 原料质量报告 ===")
	fmt.Printf("文档数:    %d\n", r.DocCount)
	fmt.Printf("总大小:    %d bytes\n", r.TotalSize)
	fmt.Printf("平均长度:  %.0f 字\n", r.AvgLength)
	fmt.Printf("有frontmatter: %d/%.0f%%\n", r.HasFrontmatter, r.HasFrontmatterPct)
	fmt.Printf("平均标题数: %.1f\n", r.AvgHeadingCount)
	fmt.Printf("平均标题深度: %.0f 级\n", r.AvgHeadingDepth)
	fmt.Printf("平均段落长: %.0f 字\n", r.AvgParaLength)
	fmt.Printf("平均列表数: %.1f\n", r.AvgListCount)
	fmt.Printf("平均表格数: %.1f\n", r.AvgTableCount)
	fmt.Printf("结构分数:  %.0f/100\n", r.StructureScore)

	verdict, reason := quality.Recommendation(r)
	fmt.Printf("建议:      %s — %s\n", verdict, reason)

	if r.DocCount > 0 && len(docs) > 0 {
		fmt.Println("\n各文档详情:")
		for _, d := range docs {
			front := " "
			if d.HasFrontmatter {
				front = "✓"
			}
			fmt.Printf("  [%s] %4d字 %2d标题 %2d列表 %2d表格  %s\n",
				front, d.Length, d.HeadingCount, d.ListCount, d.TableCount, d.Path)
		}
	}

	if *apiURL != "" {
		fmt.Println("\n--- 语义质量评估 ---")
		key := os.Getenv("DEEPSEEK_API_KEY")
		if key == "" {
			key = os.Getenv("OPENAI_API_KEY")
		}
		if key == "" {
			return fmt.Errorf("需要 DEEPSEEK_API_KEY 或 OPENAI_API_KEY")
		}
		client := llm.New(*apiURL, key, *model)

		sample := fmt.Sprintf(`文档数: %d\n平均长度: %.0f字\n结构分数: %.0f/100\n是否有frontmatter: %.0f%%\n平均标题数: %.1f\n平均列表数: %.1f`,
			r.DocCount, r.AvgLength, r.StructureScore, r.HasFrontmatterPct, r.AvgHeadingCount, r.AvgListCount)

		prompt := fmt.Sprintf(`你是一位知识工程原料评估专家。基于以下文档集的统计指标，评估这批文档作为知识工程原料的可用程度。

指标：
%s

请从 1-5 分评估三个维度：
- density（知识密度）：概念是否丰富、信息是否浓缩
- clarity（清晰度）：结构是否清晰、关系是否易提取
- suitability（适配合度）：是否适合三元组抽取

并给出结论：proceed / review / skip

输出 JSON：{"density": N, "clarity": N, "suitability": N, "conclusion": "proceed|review|skip", "reason": "简短理由"}`, sample)

		resp, err := client.Chat([]llm.Message{
			{Role: "system", Content: "你是一个严格的质量评估专家。输出 JSON。"},
			{Role: "user", Content: prompt},
		})
		if err != nil {
			return fmt.Errorf("LLM 评估失败: %w", err)
		}
		fmt.Println(resp.Choices[0].Message.Content)
	}

	return nil
}
