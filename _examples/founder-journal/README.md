# 创始人日志示例

data → info → knowl → wisdom 四层，基于创始人的日常工作日志。

## data

6 篇日志，含技术反思、工程决策、个人状态。YAML frontmatter 包含 date、mood、topics。

## info

| 文件 | 说明 |
|------|------|
| `*.md.json` | extractor 解析的结构化文档 |
| `triples.jsonl` | assessor 提取的三元组 |
| `low-confidence.jsonl` | 置信度 0.4~0.7 的暂存三元组 |

## knowl

| 文件 | 说明 |
|------|------|
| `store.jsonl` | 高置信度三元组入库后的图谱 |
| `pending-merge.jsonl` | 写入冲突待合并的三元组 |

## wisdom

| 文件 | 说明 |
|------|------|
| `faq.md` | 自动生成的 FAQ |
| `types.go` | 自动生成的 Go 类型 |
