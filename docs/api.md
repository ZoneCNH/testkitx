# API

## 模块标识

- `testkitx`：仓库名称。
- `github.com/ZoneCNH/testkitx`：Go module 路径。
- `testkitx`：公共包名。

## 公共 API

- `Config`：由用户显式提供的配置。
- `Validate`：拒绝无效配置，并返回 `ErrorKindValidation`。
- `Sanitize`：在日志或 Evidence 采集前屏蔽敏感值。
- `New`：基于显式配置创建客户端；拒绝 `nil`、canceled 和 expired context；成功时记录 `client_created_total`。
- `Close`：释放资源，并且必须幂等；成功首次关闭时记录 `client_closed_total`。
- `HealthCheck`：报告客户端健康状态，JSON 字段必须匹配 `contracts/health.schema.json`；当本次检查的 context deadline 预算短于 `Config.Timeout` 时返回 `degraded`。
- `Error`：稳定 error contract，支持 `errors.Is` / `errors.As`、`IsKind` 和 `fmt.Formatter`（`%v` / `%+v` / `%#v`）。
- `NewError` / `WrapError`：创建或包装稳定错误，包装时必须保留 cause。
- `Metrics`：注入式指标钩子；指标名必须匹配 `contracts/metrics.md`。
- `Version`：发布版本。

当前模块、下游测试采用和历史渲染 fixture 都不得依赖 `x.go`。

## 历史生成回归

`scripts/render_template.sh` 保留为历史模板兼容和 integration regression 入口。它会同步替换代码 imports、文档占位符和 module path，用于证明旧渲染链路未回归；新的下游采用应直接依赖 released `testkitx` helper，并限制在测试路径。
