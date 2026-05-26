# ROADMAP

> 从模糊的原始文档到可执行的知识产物。
>
> 每个版本交付的不是功能，是可独立验证的用户能力。

## 已完成版本

| 版本 | 层 | 交付物 |
|------|------|--------|
| v0.1.0 | 数据层 | `cmd/extractor` + `internal/document` + `internal/markdown` |
| v0.2.0 | 信息层 | `cmd/assessor` + `internal/triple/chunker/llm` |
| v0.3.0 | 知识层 | `cmd/loader` + `cmd/reasoner` + `internal/storage` + `internal/rule` |
| v0.4.0 | 智慧层 | `cmd/generator` + `internal/generator` |
| 基础设施 | — | `.github/ci.yml` `.gitignore` `Makefile` `go.mod` |

管线：`docs/*.md → cmd/extractor → cmd/assessor → cmd/loader → cmd/reasoner → cmd/generator`

---

## v0.5.0+（规划中）

### FAQ 质量评估

从数量和质量两个维度建立评估体系，驱动管线迭代。

#### 评估指标

| 维度 | 指标 | 基线 (v0.4) | 当前 (v0.5) | 目标 |
|------|------|------------|------------|------|
| 数量 | 三元组数 / 文档 | 10/6 | 60/6 | 待定 |
| 质量 | 谓语多样性 | 1 种（100%"是"） | 45 种（18%"是"） | ≥3 种 |
| 质量 | faithfulness | 4.7 | 4.4 | ≥4 |
| 质量 | usefulness | 2.4 | 2.7 | ≥4 |
| 质量 | completeness | 3.6 | 3.6 | ≥4 |

#### 改进措施（进行中）

- ✅ assessor prompt：引导提取多种谓语
- ✅ 修复默认模型名（gpt-4o-mini → deepseek-chat）
- 修复代码块解析（stripCodeBlock）
- ⬜ generator prompt：根据谓语类型调整问答模板

### 多源联合

支持文档 + 数据库表 + API 响应多源交叉建模。

- 数据源抽象接口（FileSource, DBSource, APISource）
- 联合查询与跨源关联
- 统一三元组融合策略

### 增量抽取

源文档变更后仅重抽取差异内容。

- 文件哈希缓存（SHA256 → 跳过未变更文件）
- Git-aware 差异检测（仅处理 diff 文件）
- 增量合并到已有三元组

### 自动规则挖掘

从频繁子图生成 Datalog 规则候选。

- 子图模式挖掘（Apriori / gSpan）
- 规则置信度评估
- 人工确认工作流

### 可视化看板

候选事实人工确认界面。

- Web UI（React / Vue）
- 三元组浏览器（按置信度 / 来源 / 谓词过滤）
- 批量确认 / 拒绝 / 修改

### CI 流水线

知识库变更触发自动回归推理。

- GitHub Actions：push 时运行完整 pipeline
- 推理结果快照对比（diff 检测）
- 质量门禁（置信度分布、规则覆盖率）
