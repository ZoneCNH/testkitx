# testkitx 当前状态与采用说明

> 日期：2026-06-04  
> 当前身份：L1 测试专用能力库  
> 模块路径：`github.com/ZoneCNH/testkitx`

## 1. 身份口径

`testkitx` 当前以 L1 test-only capability library 为准：它为 Go 基础库、基础设施库和仓库级 gate 提供可复用测试能力，而不是生产运行时依赖，也不是默认生成业务库的模板基座。

允许下游在以下位置使用：

- `*_test.go`
- `testkit/`
- `tools/`
- `examples/`
- CI 或测试 fixture 的临时目录

禁止下游在以下位置使用：

- 生产 `pkg/*`、`internal/*` Go 文件
- 生产二进制入口
- 运行时配置、连接池、metrics、logging 或 health 实现

若生产 import 图出现 `github.com/ZoneCNH/testkitx/pkg/testkitx` 或其子包，应视为边界失败。

## 2. 当前已具备的 L1 能力

| 能力 | 包 | 当前用途 |
|---|---|---|
| 轻量断言 | `pkg/testkitx/assertx` | 稳定失败信息、`NoError`、`ErrorIs`、`Eventually` |
| Golden 回归 | `pkg/testkitx/golden` | 默认比较，`TESTKITX_UPDATE_GOLDEN=1` 时才更新 |
| Contract 校验 | `pkg/testkitx/contract` | SHA256 contract 校验与 Evidence 写入 |
| 隔离 fixture | `pkg/testkitx/fixture` | 临时 root、HOME、module 目录和 `GOWORK=off` 环境 |
| 命令 harness | `pkg/testkitx/harness` | 超时执行命令，记录 stdout/stderr/env digest |
| 假时钟 | `pkg/testkitx/clocktest` | 确定性时间推进 |
| 可观测性 recorder | `pkg/testkitx/obstest` | 无 provider SDK 的 counters/logs 测试记录 |
| Goroutine leak 检查 | `pkg/testkitx/leaktest` | focused test 的轻量 leak 快照 |
| 生产边界扫描 | `pkg/testkitx/boundarytest` | 检查生产 Go 文件是否非法 import 测试能力 |
| Manifest fixture | `pkg/testkitx/manifesttest` | 构造、校验和写入 release manifest fixture |
| 仓库 fixture | `pkg/testkitx/repotest` | 在临时仓库中创建测试文件 |

## 3. Golden 更新规则

Golden helper 的默认行为是只读比较。只有测试进程显式设置：

```bash
TESTKITX_UPDATE_GOLDEN=1 go test ./...
```

时才允许写入或更新 golden 文件。普通 CI、release gate 和下游 smoke 不应设置该变量。这样可以避免无意中把回归输出固化为新的期望值。

## 4. Evidence 规则

L1 helper 应优先输出机器可读 Evidence，而不是只输出人工日志。当前已覆盖：

- `golden.Evidence`：记录 golden 路径、是否更新、是否匹配和实际输出 SHA256。
- `contract.Evidence`：记录 contract id、路径、SHA256、匹配状态和 metadata。
- `harness.Result`：记录命令、退出码、stdout/stderr digest、env digest、耗时和 timeout 状态。
- `manifesttest.Manifest`：记录 module、commit、gate 状态和 evidence 路径。

下游可以把这些结构写入 CI artifact、release manifest 或审计目录，但不得包含原始密钥、生产连接串或私有业务数据。

## 5. 下游采用步骤

建议按以下顺序采用：

1. 在下游仓库测试依赖中加入 `github.com/ZoneCNH/testkitx`。发布前使用 tag；本地联调可短期使用 `replace` 指向 checkout。
2. 只在 `*_test.go`、`testkit/`、`tools/` 或 `examples/` 中 import helper 包。
3. 先替换最小断言或 fixture，避免一次性重写全部测试。
4. 若使用 golden，确认普通 CI 未设置 `TESTKITX_UPDATE_GOLDEN=1`。
5. 为生产 import 图增加边界扫描；可使用 `boundarytest.ScanProductionImports` 或等价脚本。
6. 运行 `go test ./...`、`go vet ./...` 和仓库自己的 lint/typecheck gate。
7. 将 helper 产出的 Evidence 或 gate 输出作为 CI artifact 保存。
8. 在 PR 或发布说明中记录采用范围、通过的 gate 和任何未覆盖缺口。

示例测试片段：

```go
func TestDownstreamFixture(t *testing.T) {
	workspace := fixture.NewWorkspace(t, "example.test/downstream")
	assertx.Equal(t, "off", workspace.Env["GOWORK"])
}
```

## 6. 当前完成状态

本地 Standard 口径已经具备：

- L1 helper 包和对应单元测试。
- README 身份调整为测试专用能力库。
- 当前状态与采用说明文档。
- Golden opt-in 更新规则。
- Evidence 输出规则。
- 生产 import 边界约束说明。

Full 口径仍需外部事实补齐：

- 外部 CI green 的具体 workflow run URL。
- CI 上传的 release/evidence artifact URL。
- 至少一个真实下游仓库的采用 PR、commit 或测试输出证明。
- 若对外发布，还需要 tag、release manifest artifact 和变更日志条目。

这些外部事实不能由本地文档直接生成；完成声明中必须单独列明已经取得的链接或标记为未测试/待补齐。

## 7. 回归检查清单

文档或 helper 变更后至少运行：

```bash
GOWORK=off go test ./pkg/testkitx/... ./testkit/... ./contracts/...
GOWORK=off go test ./...
GOWORK=off go vet ./...
GOWORK=off make golden
git diff --check
```

如果下游采用发生变化，还应在下游仓库运行其完整测试、lint/typecheck、生产 import 扫描，并保存 CI artifact。
