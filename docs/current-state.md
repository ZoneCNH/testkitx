# 当前状态说明

## 当前实现事实

截至 2026-06-21，本仓库已经收敛为一个可编译、可测试的 Go L1 测试专用能力库，并补齐了 release readiness 文档与门禁对齐：

- `go.mod` 声明独立 module：`github.com/ZoneCNH/testkitx`。
- `pkg/testkitx` 提供公共包骨架，包括 `Config`、`Client`、`New`、`Close`、`HealthCheck`、typed error、metrics hook、option 和版本元数据。
- `internal/` 提供 validation、sanitize 和 release manifest 工具等内部实现。
- `contracts/` 锁定 config、error、health schema 和 metrics contract。
- `examples/`、`testkit/` 和 `pkg/testkitx/*` helper 子包提供 smoke 示例、测试夹具、断言、golden、contract、harness、boundary 和 Evidence 能力。
- `scripts/` 与 `Makefile` 提供 CI、boundary、contracts、security、integration、Evidence、release preflight 和 release final gate。
- `release/manifest/template.json` 是提交到源码的 manifest 模板；`release/manifest/latest.json` 是本地或 CI 生成的 Evidence artifact，不应提交。
- `CHANGELOG.md` 已记录 `v0.4.1`：同步 `FEATURES.md` 和 `ACCEPTANCE.md`、对齐 CI / Release 工作流的 `govulncheck` 版本、统一 Go 1.23 工具链默认值，以及修复临时下游模板渲染时的本地生成状态污染问题。

## 与历史 Goal 的关系

`docs/goal.md` 保留完整 Goal Prompt 和历史蓝图；当前可执行事实以 `README.md`、`docs/spec.md`、`docs/design.md`、`docs/api.md`、`docs/testing.md`、`docs/release.md`、`docs/supply-chain.md`、`contracts/`、`scripts/` 和源码为准。

如果历史蓝图与当前实现存在差异，优先相信当前代码、contracts 和 gate 输出，并把差异记录到 Evidence、review 或 retrospective，而不是直接修改历史 Goal。

## v0.4.1 变更摘要

- `FEATURES.md` 与 `ACCEPTANCE.md` 统一引用 `v0.4.1`。
- release manifest 模板版本号同步到 `v0.4.1`。
- `scripts/render_template.sh` 过滤本地 `.omc/` 与 `*.out` 生成状态，避免污染下游渲染结果。
- `CHANGELOG.md` 顶部记录本次 release readiness 同步项。

## 当前验证入口

推荐从轻到重运行：

```bash
GOWORK=off go test ./...
GOWORK=off make ci
GOWORK=off make release-check
```

正式发布前使用：

```bash
GOWORK=off make release-final-check
```

`release-final-check` 会要求 release Evidence 与当前仓库事实一致，并要求工作区为 clean。

## 已知状态约束

- 本仓库默认中文文档为主，英文标识仅保留在代码、contract、命令、包名和 gate 名称中。
- 缺少 `golangci-lint` 或 `govulncheck` 时，相关 gate 必须失败，不应被记录为跳过。
- 若运行历史模板渲染或下游采用验证，临时渲染 fixture 或采用方必须重新运行自己的 gate，并生成自己的 release Evidence 或测试 Evidence。
