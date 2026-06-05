# 测试策略

## 模块标识

- `testkitx`
- `testkitx`

## 测试策略

本仓库遵循 [测试策略母版](test-strategy.md)。当前 L1 helper 默认强制 SDD、ATDD、TDD、Contract、Boundary、Security、Integration Smoke 和 Evidence；默认增强 Property、Fuzz Smoke、Golden、Compatibility 和 Observability；Chaos、Mutation、Long Soak 和 Full E2E 只由采用方按自身 profile 启用。历史模板渲染只作为 integration regression 保留。

## L1 测试能力库规则

当前 `testkitx` 的主身份是 L1 测试专用能力库。下游只能在 `*_test.go`、`testkit/`、`tools/`、`examples/` 或 CI/fixture 临时目录中 import `github.com/ZoneCNH/testkitx/pkg/testkitx/...`；生产 Go 文件不得依赖该模块。

核心 helper 包包括：

- `assertx`：稳定断言和 eventually 检查。
- `golden`：默认只比较 golden，只有 `TESTKITX_UPDATE_GOLDEN=1` 时才更新。
- `contract`：contract SHA256 校验和机器可读 Evidence。
- `fixture`：隔离 temp root、HOME、module 目录和 `GOWORK=off` 环境。
- `harness`：命令执行与 stdout/stderr/env digest Evidence。
- `clocktest`、`obstest`、`leaktest`：确定性时间、可观测性 recorder 和 goroutine leak 检查。
- `boundarytest`、`manifesttest`、`repotest`：生产 import 边界、manifest fixture 和仓库 fixture。

采用与当前状态说明见 [当前状态与采用说明](current-state-adoption.zh-CN.md)。

## 测试模式矩阵

| 模式 | 是否默认强制 | Gate | 说明 |
|---|---:|---|---|
| SDD | 是 | `docs/spec.md` | 规格先行 |
| ATDD | 是 | `docs/testing.md` | 验收标准先行 |
| TDD / Unit | 是 | `make test` | 核心逻辑测试 |
| Race | 是 | `make race` | 并发安全 |
| Contract | 是 | `make contracts` | schema、metrics、errors |
| Boundary | 是 | `make boundary` | 模块边界 |
| Security | 是 | `make security` | `govulncheck` 和 secret scan |
| Integration Smoke | 是 | `make integration` | 历史模板渲染回归可运行 |
| Evidence | 是 | `make evidence` / `make release-check` | release manifest 与 gate 结果 |
| Property | 推荐 | `make property` | 不变量测试 |
| Fuzz Smoke | 推荐 | `make fuzz-smoke` | 边界输入测试 |
| Golden | 推荐 | `make golden` | 稳定输出回归 |
| Compatibility | 推荐 | `make contracts` | 公共契约兼容性 |
| Observability | 推荐 | `make contracts` / `make test` | metrics、health、logs |
| Chaos | 按库启用 | profile-specific | 存储和消息库 |
| Mutation | 按库启用 | critical-only | 高风险逻辑 |
| Full BDD | 不默认 | docs only | 基础库不强制 |
| Full DDD | 不作为测试模式 | boundary rule | 只保留边界思想 |

## 必需 Gate

本地执行 gate 前必须可用：

- `golangci-lint`
- `govulncheck`

缺少上述工具时，`make lint` 或 `make security` 必须失败。

- `make fmt`
- `make vet`
- `make lint`
- `make test`
- `make race`
- `make boundary`
- `make security`
- `make contracts`
- `make integration`
- `make evidence`
- `make manifest-fixture-check`

## 扩展 Gate

扩展 gate 推荐在发布前、公共 API 变更、contract 变更、schema 变更、metrics 变更和安全敏感变更时运行：

- `make property`
- `make fuzz-smoke`
- `make golden`
- `make ci-extended`
- `make release-check-extended`

`make ci` 必须保持轻量，扩展 gate 不进入默认 `make ci`。

## 必需覆盖范围

- `go test ./...` 必须覆盖公共包、`internal/`、`contracts/`、`testkit/` 和 `examples/`。
- 配置校验。
- 配置脱敏。
- typed error kind 和 wrapped cause。
- 客户端创建、取消 context、过期 context。
- 幂等关闭、zero-value client、取消 context。
- 健康与非健康状态检查。
- 健康检查 JSON 字段 contract。
- 生命周期 metrics 和健康 metrics。
- `contracts/` 与公共常量同步。
- `contracts/config.schema.json` 与 `Config` 字段映射同步。
- `scripts/render_template.sh` 生成的临时 `foundationx` 作为历史回归 fixture，可以通过 `GOWORK=off go test ./...`。
- `Config.Sanitize` 的 secret 不变量必须由 property test 覆盖。
- `Config` 边界输入必须由 fuzz-smoke 覆盖。
- `HealthStatus` JSON 公共输出必须由 golden test 锁定。
- `servicex.WaitUntil` 必须覆盖默认 deadline、调用方 deadline、已取消 context、nil ready、ready error passthrough 和 nil context。
- `contract/sql` transaction runner 必须覆盖 commit path、rollback path、operation ordering 和 tx exec error。
- `contract/timeseries` stable runner 必须覆盖 contract metric name、false stability signal 和 `Stable` error。
- `evidence` 与 `pkg/testkitx/contract` evidence writer 必须先验证再创建文件或目录，并覆盖嵌套 `WriteFile` JSON 输出。

## t.Parallel() 标准

- 无共享状态的测试函数必须调用 `t.Parallel()`。
- 表驱动测试需在 range 循环内使用 `tc := tc` 捕获循环变量后调用 `t.Parallel()`。

## 示例与 testkit Smoke

- `examples/basic` 必须输出当前 module name。
- `examples/config` 必须输出脱敏后的 secret 值。
- `examples/health` 必须输出 `healthy`。
- `testkit` 必须验证 `Config("fixture")` 生成可通过 `Validate` 的测试配置。
- `testkit.RequireNoError` 必须接受 `nil`，作为测试断言的最小契约。
- `testkit.RequireGolden` 必须比较稳定公共输出，并在 mismatch 时输出 expected 和 actual 上下文。

本仓库和临时渲染 fixture 必须保持测试独立于 `x.go`。
