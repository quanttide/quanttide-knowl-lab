# ROADMAP

> 以知识体裁（Knowledge Genre）为核心的 DIKW 管线实验室。

## 已完成版本

| 版本 | 核心能力 | 交付物 |
|------|---------|--------|
| v0.1.0 | 数据提取 | `cmd/extractor` + `internal/document` + `internal/markdown` |
| v0.2.0 | 信息评估 | `cmd/assessor` + `internal/triple` + `internal/chunker` + `internal/llm` |
| v0.3.0 | 知识入库与推理 | `cmd/loader` + `cmd/reasoner` + `internal/storage` + `internal/rule` |
| v0.4.0 | 智慧生成 | `cmd/generator` + `internal/generator` |
| v0.5.0 | 多源融合 | `cmd/federate` + `internal/source` |

v0.1–v0.5 构建了通用的 DIKW 管线。v0.6–v0.8 将管线与知识体裁对齐。

---

## v0.6.0 — 体裁识别

> 输入任意文档，自动识别其知识体裁。基于现有 `_examples/bylaw`（Manual 体裁）跑通识别管线。

核心能力：
- 体裁分类器：规则驱动（关键字 + 目录结构特征），识别 Manual / 其他
- 体裁元数据注入：在 `_examples/bylaw/` 下声明体裁标签与契约

CLI 入口：
```
go run ./cmd/genre detect -input _examples/bylaw
```

交付物：
- `cmd/genre/` — 体裁 CLI（子命令模式：`detect` / `validate` / `transcode`）
- `internal/genre/` — 体裁模型、分类器
- `_examples/bylaw/genre.json` — Manual 体裁元数据
- `cmd/genre/main_test.go` — 各子命令集成测试

---

## v0.7.0 — 体裁校验

> 校验文档是否符合其体裁契约。以 Manual 体裁的契约（形式、粒度、生命周期阶段）验证 bylaw 的结构完整性。

核心能力：
- 结构校验：data/info/knowl/wisdom 各层齐全
- 字段校验：每层关键字段存在
- 跨层一致性校验：info 问答对 ↔ knowl 三元组对照

CLI 入口：
```
go run ./cmd/genre validate -input _examples/bylaw -schema _examples/bylaw/genre.json
```

交付物：
- `cmd/genre validate` — 校验子命令
- `internal/genre/validator/` — 校验引擎
- Manual 体裁契约定义（内嵌于 `internal/genre/`）
- 校验报告格式（通过 / 警告 / 错误）

---

## v0.8.0 — 体裁转换

> 将 Manual 体裁的知识转换为 Archive 体裁的结构化 Instance。验证转换后产物符合 Archive 契约。

核心能力：
- Manual → Archive 单向转换（bylaw → 结构化 Instance）
- 转换后 Archive 的体裁自校验

CLI 入口：
```
go run ./cmd/genre transcode -from manual -to archive -input _examples/bylaw -output _examples/bylaw/archive
```

交付物：
- `cmd/genre transcode` — 转换子命令
- `internal/genre/transcoder/` — 转换器
- Archive 体裁契约定义
- 收集第二个体裁的原始数据（如 `_examples/journal/`），为多体裁验证提供基础

---

## 路线总览

```
版本    能力          输入 → 输出                    体裁维度
────────────────────────────────────────────────────────────
v0.1   数据提取      docs/*.md → parsed              无
v0.2   信息评估      parsed → triples                无
v0.3   知识推理      triples → store                 无
v0.4   智慧生成      store → code/faq/api            无
v0.5   多源融合      multi-source → unified triples  无
────────────────────────────────────────────────────────────
v0.6   体裁识别      _examples/* → 体裁标签          识别
v0.7   体裁校验      体裁文档 + schema → 校验报告    结构化验证
v0.8   体裁转换      体裁A → 体裁B                   跨体裁转换
```
