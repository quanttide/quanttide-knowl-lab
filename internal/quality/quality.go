package quality

import (
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Report struct {
	DocCount           int
	TotalSize          int64
	AvgLength          float64
	HasFrontmatter     int
	HasFrontmatterPct  float64
	AvgHeadingDepth    float64
	AvgHeadingCount    float64
	AvgParaLength      float64
	AvgListCount       float64
	AvgTableCount      float64
	StructureScore     float64
}

type DocMetrics struct {
	Path          string
	Size          int64
	Length        int
	HasFrontmatter bool
	MaxHeading    int
	HeadingCount  int
	ParaCount     int
	AvgParaLen    float64
	ListCount     int
	TableCount    int
}

func AnalyzeDir(dir string) (*Report, []DocMetrics, error) {
	var docs []DocMetrics
	totalSize := int64(0)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".md" && ext != ".txt" {
			return nil
		}
		m, err := analyzeFile(path)
		if err != nil {
			return nil
		}
		docs = append(docs, *m)
		totalSize += m.Size
		return nil
	})

	if err != nil {
		return nil, nil, err
	}
	if len(docs) == 0 {
		return &Report{DocCount: 0}, nil, nil
	}

	r := &Report{DocCount: len(docs), TotalSize: totalSize}
	for _, d := range docs {
		r.AvgLength += float64(d.Length)
		r.AvgHeadingDepth += float64(d.MaxHeading)
		r.AvgHeadingCount += float64(d.HeadingCount)
		r.AvgParaLength += d.AvgParaLen
		r.AvgListCount += float64(d.ListCount)
		r.AvgTableCount += float64(d.TableCount)
		if d.HasFrontmatter {
			r.HasFrontmatter++
		}
	}
	n := float64(len(docs))
	r.AvgLength /= n
	r.AvgHeadingDepth /= n
	r.AvgHeadingCount /= n
	r.AvgParaLength /= n
	r.AvgListCount /= n
	r.AvgTableCount /= n
	r.HasFrontmatterPct = float64(r.HasFrontmatter) / n * 100

	// StructureScore: 0-100 composite
	// Weighted: has_frontmatter(20) + heading_count(20) + para_length(20) + lists(20) + tables(20)
	if n == 0 {
		return r, docs, nil
	}
	score := 0.0
	if r.HasFrontmatterPct > 50 {
		score += 20
	} else if r.HasFrontmatterPct > 0 {
		score += 10
	}
	if r.AvgHeadingCount >= 3 {
		score += 20
	} else if r.AvgHeadingCount >= 1 {
		score += 10
	}
	if r.AvgParaLength >= 100 && r.AvgParaLength <= 2000 {
		score += 20
	} else if r.AvgParaLength > 0 {
		score += 10
	}
	if r.AvgListCount >= 1 {
		score += 20
	} else if r.AvgListCount > 0 {
		score += 10
	}
	if r.AvgTableCount >= 1 {
		score += 20
	} else if r.AvgTableCount > 0 {
		score += 10
	}
	r.StructureScore = score

	return r, docs, nil
}

func analyzeFile(path string) (*DocMetrics, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	text := string(data)
	m := &DocMetrics{Path: path, Size: int64(len(data))}
	m.Length = utf8.RuneCountInString(text)

	lines := strings.Split(text, "\n")
	m.HasFrontmatter = strings.HasPrefix(strings.TrimSpace(text), "---")

	var paras []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			m.HeadingCount++
			depth := 0
			for _, r := range line {
				if r == '#' {
					depth++
				} else {
					break
				}
			}
			if depth > m.MaxHeading {
				m.MaxHeading = depth
			}
			continue
		}
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*") || (len(line) > 1 && line[0] >= '0' && line[0] <= '9' && line[1] == '.') {
			m.ListCount++
			continue
		}
		if strings.HasPrefix(line, "|") {
			m.TableCount++
			continue
		}
		paras = append(paras, line)
	}
	if len(paras) > 0 {
		m.ParaCount = len(paras)
		totalParaLen := 0
		for _, p := range paras {
			totalParaLen += utf8.RuneCountInString(p)
		}
		m.AvgParaLen = float64(totalParaLen) / float64(len(paras))
	}
	return m, nil
}

func Recommendation(r *Report) (string, string) {
	if r.DocCount == 0 {
		return "skip", "无文档"
	}
	if r.DocCount < 3 {
		return "review", "文档数量少（<3），可能不足以抽取有意义的知识"
	}
	if r.StructureScore < 30 {
		return "review", "结构分数低，文档可能缺乏清晰层次"
	}
	if r.AvgLength < 200 {
		return "review", "文档过短（平均 < 200 字），内容密度可能不足"
	}
	if r.AvgHeadingCount < 1 && r.HasFrontmatterPct < 50 {
		return "review", "既无标题结构也无 frontmatter，原料质量存疑"
	}
	return "proceed", "原料可用"
}
