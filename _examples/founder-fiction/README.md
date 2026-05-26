# 创始人小说示例

data → info → knowl → wisdom 四层，基于创始人的职场言情小说章节。

## data

13 篇小说章节，含对话、叙事、情感描写。与 founder-journal 的技术反思不同，
fiction 提供创意写作领域的语料，验证知识抽取在不同文体上的表现。

## info

| 文件 | 说明 |
|------|------|
| `*.md.json` | extractor 解析的结构化文档 |
| `triples.jsonl` | assessor 提取的三元组 |
| `low-confidence.jsonl` | 置信度 0.4~0.7 的暂存三元组 |
| `errors.jsonl` | assessor 处理失败的记录（LLM 超时、格式异常等） |

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
