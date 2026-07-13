# ROADMAP

> 以知识体裁（Knowledge Genre）为核心的 DIKW 管线实验室。
>
> 每个版本交付的不是功能，是可独立验证的体裁转换能力。

## 已完成版本

| 版本 | 核心能力 | 交付物 |
|------|---------|--------|
| v0.1.0 | 数据提取 | `cmd/extractor` + `internal/document` + `internal/markdown` |
| v0.2.0 | 信息评估 | `cmd/assessor` + `internal/triple` + `internal/chunker` + `internal/llm` |
| v0.3.0 | 知识入库与推理 | `cmd/loader` + `cmd/reasoner` + `internal/storage` + `internal/rule` |
| v0.4.0 | 智慧生成 | `cmd/generator` + `internal/generator` |
| v0.5.0 | 多源融合 | `cmd/federate` + `internal/source` |

---

## v0.6.0 — 体裁识别（Genre Detection）

> 输入任意文档，自动识别其知识体裁。

核心能力：
- 体裁分类器：判断输入属于语境/日志/档案/手册或未分类
- 体裁特征提取：从文档结构、措辞模式、元数据中提取体裁特征
- 混合体裁检测：单一文档可能包含多种体裁（如语境中嵌入日志引用）

CLI 入口：
```
go run ./cmd/genre detect -input _examples/founder-journal/data
```

交付物：
- `cmd/genre/` — 体裁检测 CLI
- `internal/genre/` — 体裁模型、分类器、特征提取器
- 体裁分类准确率 ≥ 80%（基于现有 _examples 场景）

---

## v0.7.0 — 体裁校验（Genre Validation）

> 校验任意体裁的知识是否符合 Domain/Ontology/Instance 契约。

核心能力：
- 结构校验：字段完整、引用有效、类型匹配
- 跨体裁一致性校验：同一知识在语境和档案中的描述不自相矛盾
- 本体契约检查：Instance 引用的 Ontology 存在且字段对齐

CLI 入口：
```
go run ./cmd/genre validate -input data/journal/2026-07-13.md -schema knowledge-schema.json
```

交付物：
- `cmd/genre validate` — 校验子命令
- `internal/genre/validator/` — 校验引擎
- 知识体裁的 JSON Schema 定义（三要素：形式、粒度、生命周期阶段）

---

## v0.8.0 — 体裁转换（Genre Transcoding）

> 同一条知识，在不同体裁之间双向转换。

核心能力：
- 语境 → 档案：抽取叙述中的设计决策，结构化为 Instance
- 日志 → 档案：按时序聚合日志中的判断，沉淀为 Instance
- 档案 → 手册：将结构化知识编译为可执行的指令序列
- 统一中间表示：所有转换经过 Domain/Ontology/Instance 骨架

CLI 入口：
```
go run ./cmd/genre transcode -from journal -to archive -input journal/entry.md -output archive/entry.json
```

交付物：
- `cmd/genre transcode` — 体裁转换子命令
- `internal/genre/transcoder/` — 转换器（每个体裁对一条转换路径）
- 基准场景：founder-journal（日志 → 档案）完整管线演示

---

## v0.9.0 — 体裁追溯（Genre Trace）

> 从档案中的一个 Instance，追溯到它在语境和日志中的原始形态。

核心能力：
- 向上追溯：Instance → 日志条目 → 语境讨论（来源回溯）
- 向下追溯：语境讨论 → 日志决策 → 档案 Instance（沉淀追踪）
- 变更历史：同一知识在不同体裁中的演化过程

CLI 入口：
```
go run ./cmd/genre trace -id <instance-uuid> -depth all
```

交付物：
- `cmd/genre trace` — 追溯子命令
- 追溯链可视化输出（Mermaid / DOT 格式）
- 端到端演示：founder-fiction（创意写作领域全体裁追溯）

---

## 路线总览

```
版本    能力          输入 → 输出                    体裁维度
────────────────────────────────────────────────────────────
v0.1   数据提取      docs/*.md → parsed              无（原材料）
v0.2   信息评估      parsed → triples                无（原材料）
v0.3   知识推理      triples → store                 无（原材料）
v0.4   智慧生成      store → code/faq/api            无（原材料）
v0.5   多源融合      multi-source → unified triples  无（原材料）
────────────────────────────────────────────────────────────
v0.6   体裁识别      任意文档 → 体裁标签             识别体裁
v0.7   体裁校验      体裁文档 → 校验报告             结构化验证
v0.8   体裁转换      体裁A → 体裁B                   跨体裁转换
v0.9   体裁追溯      Instance → 全体裁链路           同源追溯
```

v0.1~v0.5 构建了通用的 DIKW 管线。v0.6~v0.9 将管线与知识体裁对齐——让管线不仅能加工知识，还能理解知识的呈现形态，并在不同形态之间自由转换。
