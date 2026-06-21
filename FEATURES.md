# FEATURES

`testkitx` 是 `github.com/ZoneCNH/testkitx` 的 L1 测试专用能力库。当前发布版本以 `pkg/testkitx/version.go`、`.repo-contract.yaml` 和 `CHANGELOG.md` 为准，彼此应保持同一版本号。

- 当前版本：`v0.4.1`

## 当前特性

- 断言辅助：`assertx`、`Eventually`、稳定失败信息与 helper 栈保留。
- Golden 回归：bytes / JSON 比较，默认只读，显式允许更新。
- Contract / Evidence：SHA256 校验、manifest 生成和 release evidence 写入。
- Fixture：隔离 workspace、`HOME`、module 目录和 `GOWORK=off` 环境。
- Harness：命令执行、stdout / stderr / env digest Evidence。
- 可观测性：fake clock、metrics recorder、goroutine leak 检查。
- 边界与仓库辅助：生产 import 边界扫描、manifest fixture、repo fixture。

## 设计边界

- 仅供测试、工具、示例和 CI / fixture 使用。
- 不作为生产运行时依赖。
- 不要求下游从模板生成代码。
