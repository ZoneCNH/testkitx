# 当前状态说明

## 当前实现事实

截至 2026-06-04，本仓库已经收敛为一个可编译、可测试的 Go 基础库模板：

- `go.mod` 声明独立 module：`github.com/ZoneCNH/testkitx`。
- `pkg/testkitx` 提供公共包骨架，包括 `Config`、`Client`、`New`、`Close`、`HealthCheck`、typed error、metrics hook、option 和版本元数据。
- `internal/` 提供 validation、sanitize 和 release manifest 工具等内部实现。
- `contracts/` 锁定 config、error、health schema 和 metrics contract。
- `examples/` 与 `testkit/` 提供 smoke 示例、测试夹具和断言。
- `scripts/` 与 `Makefile` 提供 CI、boundary、contracts、security、integration、Evidence 和 release preflight gate。
- `release/manifest/template.json` 是提交到源码的 manifest 模板；`release/manifest/latest.json` 是本地或 CI 生成的 Evidence artifact，不应提交。
- `CHANGELOG.md` 已记录 `v0.3.0`：`fixture.WriteOrFatal(t)`、`Error` 的 `fmt.Formatter` 支持、`t.Parallel()` 测试并行化、CI 覆盖率报告、`releasemanifest` 文件拆分和 README badges 已纳入当前治理状态。

## 与历史 Goal 的关系

`docs/goal.md` 保留完整 Goal Prompt 和历史蓝图；当前可执行事实以 `README.md`、`docs/spec.md`、`docs/design.md`、`docs/api.md`、`docs/testing.md`、`docs/release.md`、`docs/supply-chain.md`、`contracts/`、`scripts/` 和源码为准。

如果历史蓝图与当前实现存在差异，优先相信当前代码、contracts 和 gate 输出，并把差异记录到 Evidence、review 或 retrospective，而不是直接修改历史 Goal。

## v0.3.0 变更摘要

- `fixture.WriteOrFatal(t)` 为推荐的 fixture 写入 API，替代 panic 行为。
- `Error` 支持 `fmt.Formatter`（`%v` / `%+v` / `%#v`）。
- 测试函数使用 `t.Parallel()` 并行执行。
- CI 生成覆盖率报告（`go test -coverprofile`）。
- `releasemanifest/main.go`（532 行）拆分为 6 个文件（最大 131 行）。
- README 添加 CI/Security/ReportCard/License/Version badges。

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
- 生成下游库后，下游库必须重新运行自己的 gate 并生成自己的 release Evidence。
