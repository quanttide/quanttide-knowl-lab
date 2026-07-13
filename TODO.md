# TODO

> 当前阶段：bylaw 是唯一示例场景（体裁：Manual），以此为起点跑通体裁管线。

## v0.6.0 — 体裁识别

- [ ] `internal/genre/` 包骨架
  - `genre.go` — 体裁模型（Genre 类型、特征接口）
  - `classifier.go` — 规则驱动分类器（关键词 + 目录结构特征）
  - `genre_test.go` — 单元测试（含 bylaw → Manual 识别用例）
- [ ] `_examples/bylaw/genre.json` — 体裁元数据
  - 定义 Manual 体裁的契约：形式、粒度、生命周期阶段
- [ ] `cmd/genre/` 骨架（子命令模式）
  - `main.go` — 子命令分发
  - `detect.go` — `detect` 子命令 + `run()` 函数
  - `main_test.go` — 集成测试（至少能正确识别 bylaw → Manual）
- [ ] 更新 `Makefile` — 添加 `genre` target

## v0.7.0 — 体裁校验

- [ ] `internal/genre/validator/` 校验引擎
  - `validator.go` — 校验器（接收 Genre + 文档路径 → 校验报告）
  - `rules.go` — Manual 体裁的具体校验规则（结构完整、字段完整、跨层一致）
  - `validator_test.go` — 单元测试
- [ ] `cmd/genre validate.go` — `validate` 子命令 + `run()` 函数
  - 输入：体裁文档目录
  - 输出：校验报告（通过 / 警告 / 错误）
- [ ] 端到端验证：对 `_examples/bylaw` 运行 validate，确认通过
- [ ] `cmd/genre/main_test.go` 增加 validate 集成测试

## v0.8.0 — 体裁转换

- [ ] 设计并定义 Archive 体裁契约
  - 结构化 Instance 的字段定义
  - Archive 与 Manual 的对应关系
- [ ] `internal/genre/transcoder/manual_to_archive.go` — Manual → Archive 转换器
- [ ] `cmd/genre transcode.go` — `transcode` 子命令 + `run()` 函数
- [ ] 验证转换后的 Archive 符合 Archive 体裁契约
- [ ] `cmd/genre/main_test.go` 增加 transcode 集成测试
- [ ] 收集第二个体裁的原始数据
  - 在 `_examples/` 下新增 `journal/` 或从 `docs/context` 引入实例
  - 多个体裁才能验证识别、校验、转换的通用性

## 文档与基础设施

- [ ] 更新 `AGENTS.md` — 添加 `cmd/genre` 和 `internal/genre/` 目录说明
- [ ] 更新 `README.md` — 添加体裁管线说明与 v0.6–v0.8 快速开始示例
- [ ] 更新 `Makefile` — 为 v0.7、v0.8 添加对应 target
