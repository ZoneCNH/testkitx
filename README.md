# testkitx

[![CI](https://github.com/ZoneCNH/testkitx/actions/workflows/ci.yml/badge.svg)](https://github.com/ZoneCNH/testkitx/actions/workflows/ci.yml)
[![Security](https://github.com/ZoneCNH/testkitx/actions/workflows/security.yml/badge.svg)](https://github.com/ZoneCNH/testkitx/actions/workflows/security.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ZoneCNH/testkitx)](https://goreportcard.com/report/github.com/ZoneCNH/testkitx)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ZoneCNH/testkitx)](go.mod)

`testkitx` 是 `github.com/ZoneCNH/testkitx` 的 **L1 测试专用能力库**。它为 Go 基础库和基础设施库提供可复用的测试断言、golden 回归、contract hash、隔离 fixture、命令 harness、fake clock、可观测性 recorder、goroutine leak 检查、生产 import 边界扫描、manifest fixture 和仓库文件 fixture。

`testkitx` 不是生产运行时基础库，也不是用于生成业务库的默认模板基座。下游库只能在测试代码、测试工具、示例或显式测试夹具中使用它；生产包不得 import `github.com/ZoneCNH/testkitx/pkg/testkitx/...`。

## 目标

本模块的目标是让基础库从第一天就具备可验证的 L1 测试语义：

- 断言失败信息稳定、可读，并保留 `testing.TB.Helper()` 调用栈。
- Golden 文件默认只比较；只有显式设置 `TESTKITX_UPDATE_GOLDEN=1` 时才允许更新。
- Contract、harness、golden 和 manifest 辅助能力返回机器可读 Evidence，便于 CI artifact、release manifest 和审计串联。
- Fixture 为每个测试创建隔离的临时目录、`HOME`、module 目录和 `GOWORK=off` 环境。
- Boundary helper 可扫描下游仓库，防止测试能力泄漏到生产 import 图。
- 所有测试能力保持独立，不依赖 `x.go`、真实凭据或生产连接。

## 非目标

- 不作为生产运行时依赖。
- 不自动读取生产密钥或环境。
- 不创建隐藏全局客户端。
- 不承载 x.go 业务模型。
- 不替代下游库自己的 L2/L3 integration、chaos、soak 或真实外部系统测试。
- 不在未设置 `TESTKITX_UPDATE_GOLDEN=1` 时改写 golden 文件。

## 能力包

公共 L1 测试能力位于 `pkg/testkitx/`：

- `assertx`：轻量断言、`NoError`、`ErrorIs` 和 `Eventually`。
- `golden`：bytes/JSON golden 比较与 opt-in 更新，输出 hash Evidence。
- `contract`：contract 文件 SHA256 校验与 Evidence 写入。
- `fixture`：隔离 workspace、HOME、module 和环境变量。
- `harness`：带超时的命令执行与 stdout/stderr/env digest Evidence。
- `clocktest`：确定性 fake clock。
- `obstest`：无 provider SDK 的 counters/log recorder。
- `leaktest`：轻量 goroutine leak 快照与校验。
- `boundarytest`：扫描生产 Go 文件中的非法测试库 import。
- `manifesttest`：构造、写入和读取 release manifest fixture。
- `repotest`：仓库 fixture 文件写入辅助。

仓库仍保留 `Config`、`Client`、`HealthCheck`、metrics、contracts、release manifest 和模板文档等历史模板资产；它们用于兼容当前 gate 和回归基线。新的文档身份以 L1 测试专用能力库为准。

## 使用边界

允许使用位置：

- `*_test.go`
- `testkit/`
- `tools/`
- `examples/`
- 专门用于测试或 CI 的临时 fixture 目录

禁止使用位置：

- `pkg/*`、`internal/*` 等生产 Go 文件
- 生产二进制入口
- 运行时配置、连接、metrics 或日志实现

下游仓库可使用 `boundarytest.ScanProductionImports` 或自有脚本扫描生产 import 图。发现生产文件 import `github.com/ZoneCNH/testkitx/pkg/testkitx` 或其子包时，应视为失败。

## 文档入口

- [身份说明](docs/identity.md)：项目身份、模板定位和不可越界的生产约束。
- [当前状态](docs/current-state.md)：当前实现事实、历史 Goal 关系和验证入口。
- [采用说明](docs/adoption.md)：下游基础库采用路径、检查清单和风险。
- [规格](docs/spec.md)：模板能力、验收标准和可追踪性。
- [设计](docs/design.md)：模块边界、公共 API、错误、健康检查和指标设计。
- [API](docs/api.md)：`Config`、`Client`、typed error、health JSON 和 metrics contract。
- [配置](docs/config.md)：显式配置、validation 和脱敏规则。
- [生成](docs/generation.md)：从模板渲染 `foundationx` 等具体基础库。
- [错误模型](docs/errors.md)：`ErrorKind`、`NewError`、`WrapError` 和重试语义。
- [可观测性](docs/observability.md)：指标名、健康状态和 JSON 字段。
- [测试策略母版](docs/test-strategy.md)：Required、Extended 和 profile-specific gates。
- [规格](docs/spec.md)：历史模板能力、验收标准和可追踪性。
- [设计](docs/design.md)：历史模板边界、公共 API、错误、健康检查和指标设计。
- [API](docs/api.md)：历史模板 API 与 contract 说明。
- [配置](docs/config.md)：显式配置、validation 和脱敏规则。
- [生成](docs/generation.md)：历史模板渲染说明；当前 L1 采用不要求下游从模板生成。
- [供应链](docs/supply-chain.md)：可校验 release Evidence、源码摘要、contract 指纹、依赖清单和 CI artifact。
- [发布](docs/release.md)：`release-check`、manifest 字段和 Evidence 规则。

## 命令

本地完整 gate 需要安装 `golangci-lint` 和 `govulncheck`；CI 会显式安装这两个工具。缺少任一工具时，`make lint` 或 `make security` 必须失败，不允许把必需 gate 记录为跳过。

```bash
GOWORK=off go test ./...
GOWORK=off go vet ./...
GOWORK=off make golden
GOWORK=off make ci
GOWORK=off make release-check
```

只验证 L1 helper 包时可运行：

```bash
GOWORK=off go test ./pkg/testkitx/... ./testkit/... ./contracts/...
```

## 采用示例

下游测试可直接导入具体 helper：

```go
package example_test

import (
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/assertx"
	"github.com/ZoneCNH/testkitx/pkg/testkitx/fixture"
)

func TestConfigFixture(t *testing.T) {
	workspace := fixture.NewWorkspace(t, "example.test/downstream")
	assertx.Equal(t, "off", workspace.Env["GOWORK"])
}
```

发布 tag 前推荐使用 released 版本；在本仓库发布前或本地联调时，下游可临时使用 `replace` 指向本地 checkout，并在发布前移除。

## Evidence

本仓库的 L1 完成声明需要同时给出本地 gate 和可审计 Evidence：

- `GOWORK=off go test ./...` 覆盖 helper、contracts、examples 和历史模板基线。
- `GOWORK=off go vet ./...` 覆盖 Go 静态检查。
- `GOWORK=off make golden` 锁定稳定输出且不默认更新 golden。
- Contract、harness、golden、manifest helper 通过结构化 Evidence 类型承载审计字段。
- Full 口径仍需要外部 CI artifact URL、下游真实仓库采用证明和发布 tag/manifest artifact；这些不是本地文档变更能单独证明的事实。

最终对外完成声明应包含 `DONE with evidence:`，并列出本地 gate、外部 CI artifact、下游采用证据和任何未完成缺口。
