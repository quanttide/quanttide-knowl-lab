# TODO

> 当前阶段：bylaw 是唯一示例场景（体裁：手册），以此为起点跑通体裁概念。

## v0.6 — 体裁识别

- [ ] 给 bylaw 打上体裁标签（Manual）
  - 在 `_examples/bylaw/` 下添加体裁元数据（`genre.json` 或 frontmatter）
  - 定义 Manual 体裁的契约：形式、粒度、生命周期阶段
- [ ] `cmd/genre detect` 最小实现
  - 输入：任意文档路径
  - 输出：体裁标签 + 置信度
  - 至少能正确识别 bylaw（Manual）
- [ ] `internal/genre/` 包骨架
  - `genre.go` — 体裁模型定义
  - `classifier.go` — 简单分类器（规则驱动，先用关键字/结构特征）

## v0.7 — 体裁校验

- [ ] 验证 bylaw 是否符合 Manual 体裁的契约
  - 结构完整性（data/info/knowl/wisdom 各层齐全）
  - 字段完整性（每层的关键字段是否存在）
  - 跨层一致性（info 中的问答对是否与 knowl 中的三元组对应）
- [ ] `cmd/genre validate` 最小实现
  - 输入：体裁文档目录 + Schema
  - 输出：校验报告（通过/警告/错误）

## v0.8 — 体裁转换

- [ ] 设计由 Manual 到 Archive 的转换路径（bylaw → 结构化 Instance）
- [ ] `cmd/genre transcode` 骨架
- [ ] 验证转换后的 Archive 是否符合体裁契约

## 数据建设

- [ ] 收集第二个体裁的数据（如 journal 或 context 中的实例）
  - 多个体裁才能验证识别、校验和转换的通用性

## 基础设施

- [ ] 更新 `AGENTS.md`，反映体裁能力
- [ ] 更新 `README.md`，添加体裁管线说明
