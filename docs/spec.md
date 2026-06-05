# SPEC-testkitx-v1.0

## 需求

- 为 L1 测试专用能力库提供独立 Go module。
- 提供 `Config`、`Validate`、`Sanitize`、`Client`、`New`、`Option`、`HealthCheck`、错误模型、指标钩子和版本元数据。
- 提供可在下游测试路径采用的断言、fixture、golden、contract、harness、clock、observability、leak、boundary、manifest 和 repo helper。
- `Validate`、`New`、`Close` 和 `HealthCheck` 必须返回或记录可分类的生产语义，包括 typed error、幂等关闭、上下文取消和健康状态。
- 提供 Harness Gate 脚本、历史模板回归脚本、CI 工作流、contracts、examples、Evidence artifact、release 和复盘文档。

## 验收标准

- `GOWORK=off go test ./...` 和 `GOWORK=off go test -race ./...` 通过。
- `GOWORK=off make release-check` 通过，并以 `CHECK_STATUS=passed` 生成未提交的 `release/manifest/latest.json` 与 `release/manifest/latest.json.sha256` Evidence artifact。
- `contracts/config.schema.json` 与 `Config` 字段映射保持一致，`timeout_ms` 映射到 `Config.Timeout`。
- `contracts/error.schema.json`、`contracts/health.schema.json` 和 `contracts/metrics.md` 与公共常量保持一致。
- `scripts/render_template.sh` 作为历史 regression 入口，可以生成 `foundationx` 形态并通过 `GOWORK=off go test ./...`。
- 模块不得依赖 `github.com/bytechainx/x.go` 或 `github.com/ZoneCNH/x.go`。
- 模块不得隐式读取生产密钥。
- 下游生产 import 图不得依赖 `github.com/ZoneCNH/testkitx/pkg/testkitx` 或其子包。

## 非目标

- 不包含业务模型、生产连接默认值和隐藏全局客户端。

## 可追踪性

- 目标：`GOAL-20260601-001`
- 模板占位符：`testkitx`、`github.com/ZoneCNH/testkitx`、`testkitx`
