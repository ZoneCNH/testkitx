# 身份说明

## 项目身份

`testkitx` 是 `github.com/ZoneCNH/testkitx` 的 L1 测试专用能力库。它为 Go 基础库和基础设施库提供可复用的测试断言、golden 回归、contract hash、隔离 fixture、命令 harness、fake clock、可观测性 recorder、goroutine leak 检查、生产 import 边界扫描、manifest fixture 和仓库文件 fixture。

`testkitx` 不再作为下游生产库的默认生成模板，也不承担生产运行时基础库职责。下游只能在测试代码、测试工具、示例或显式测试夹具中使用它；生产包不得 import `github.com/ZoneCNH/testkitx/pkg/testkitx/...`。

## 与 xlib-standard 的关系

`xlib-standard` 是基础库体系的 Standard Source。`testkitx` 需要同步其中适用于 L1 测试库的共享标准，例如 CI/security gate、Docker toolchain contract、release Evidence 口径、文档语言规则和边界约束。

同步范围必须受 L1 身份限制：可以同步标准、契约和验证入口；不得同步会把 `testkitx` 重新变成生产模板、业务库、真实 runtime adapter、隐式密钥读取器或 `x.go` 依赖的能力。

## 历史模板资产

仓库仍保留 `Config`、`Client`、`HealthCheck`、metrics、contracts、release manifest 和 `scripts/render_template.sh` 等历史模板资产。这些资产当前用于兼容已有 gate、integration regression 和迁移期基线，不代表新的采用方式。

新增能力和对外说明应优先服务 L1 测试专用身份。若历史模板资产与 L1 身份冲突，应先隔离为 regression fixture 或迁移计划，再考虑删除或改造，避免破坏现有 gate。

## 边界

`testkitx` 不负责下游库的生产连接、密钥读取、业务语义或 runtime client 封装。任何派生库、下游库或测试夹具都必须由调用方显式传入配置；本仓库不得自动读取 `/home/k8s/secrets/env/*`，也不得写入真实密钥、生产地址或隐藏全局 client。

禁止事项：

- 不 import `github.com/bytechainx/x.go`、`github.com/ZoneCNH/x.go` 或任何 `x.go/internal/*` 包。
- 不在生产 Go 文件中依赖 `pkg/testkitx` 测试 helper。
- 不把测试 fake、recorder、fixture 或 harness 注册成生产默认实现。
- 不以 README、计划文档或本地 dry-run 替代可复验 Evidence。

## 完成身份

当本仓库声明完成时，声明依据必须是可复验的 Evidence，而不是口头状态。完成语句沿用仓库约定：`DONE with evidence:`，并附带测试、contract、boundary、security、CI、Docker toolchain、下游 test-only 采用证明和 release manifest 证据。
