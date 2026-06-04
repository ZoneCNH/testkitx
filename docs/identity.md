# 身份说明

## 项目身份

`testkitx` 是 x.go 基础库体系的独立 Go 模板仓库，用来沉淀可复用的基础库骨架、公共 API 约定、Harness Gate、release Evidence 和复盘工件。它本身不是业务库，也不是某个具体基础设施适配器。

## 模板定位

`testkitx` 的核心身份是“生产级共享基础库基座”：

- 为 `foundationx`、`configx`、`observex`、`postgresx`、`kafkax`、`redisx`、`taosx`、`ossx` 等后续库提供统一生成源。
- 锁定基础库必须继承的最小生产语义，包括显式 `Config`、稳定 validation error、typed error、幂等 `Close`、`HealthCheck`、metrics contract 和 release manifest。
- 保持独立 Go module，不 import `github.com/bytechainx/x.go` 或 `github.com/ZoneCNH/x.go`，不包含 `x.go` 业务模型。

## 边界

`testkitx` 只负责模板层能力，不负责下游库的生产连接、密钥读取或业务语义。派生库必须由调用方显式传入配置；模板不得自动读取 `/home/k8s/secrets/env/*`，也不得写入真实密钥、生产地址或隐藏全局 client。

## 完成身份

当本仓库声明完成时，声明依据必须是可复验的 Evidence，而不是口头状态。完成语句沿用仓库约定：`DONE with evidence:`，并附带测试、contract、boundary、security、integration 和 release manifest 证据。
