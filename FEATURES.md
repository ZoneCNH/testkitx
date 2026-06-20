# FEATURES

`testkitx` 是 `github.com/ZoneCNH/testkitx` 的 **L1 测试专用能力库**。本文件记录当前已审计的功能边界，作为 release readiness 的功能清单。

## 已确认的功能面

- `pkg/testkitx/assertx`：轻量断言、`NoError`、`ErrorIs`、`Eventually` 与格式化错误输出。
- `pkg/testkitx/golden`：bytes / JSON golden 比较、按需更新与 hash Evidence。
- `pkg/testkitx/contract`：contract 文件 SHA256 校验与结构化 Evidence 写入。
- `pkg/testkitx/fixture`：隔离 workspace、HOME、module 和环境变量的测试夹具。
- `pkg/testkitx/harness`：命令执行、超时控制与 stdout / stderr / env digest Evidence。
- `pkg/testkitx/clocktest`：确定性 fake clock。
- `pkg/testkitx/obstest`：无需 provider SDK 的 counters / log recorder。
- `pkg/testkitx/leaktest`：goroutine leak 快照与校验。
- `pkg/testkitx/boundarytest`：扫描生产 Go 文件中的非法测试库 import。
- `pkg/testkitx/manifesttest`：release manifest 生成、写入和读取夹具。
- `pkg/testkitx/repotest`：仓库 fixture 文件写入辅助。

## 已确认的约束

- 仅用于测试、fixture、工具和示例，不作为生产运行时依赖。
- 不应隐式读取生产密钥路径，也不应隐藏创建全局客户端。
- golden 默认只比较，只有显式设置 `TESTKITX_UPDATE_GOLDEN=1` 才允许更新。
- contract、golden、harness 和 manifest helper 应继续返回机器可读 Evidence。
- 生产包不得 import `github.com/ZoneCNH/testkitx/pkg/testkitx/...`。

## Release readiness 相关入口

- `README.md`
- `docs/release.md`
- `docs/current-state.md`
- `docs/spec.md`
- `docs/test-strategy.md`

