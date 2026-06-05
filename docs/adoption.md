# 采用说明

## 适用场景

`testkitx` 适用于下游 Go 基础库和基础设施库的 L1 测试代码。典型用途包括：

- 在 `*_test.go` 中复用 `assertx`、`golden`、`contract`、`fixture`、`harness`、`clocktest`、`obstest`、`leaktest`、`manifesttest` 和 `repotest`。
- 在 `testkit/` 或测试专用工具中构造隔离 workspace、module fixture、manifest fixture 或命令执行 Evidence。
- 在 CI 中使用 `boundarytest` 或等价脚本扫描生产 import 图，防止测试库泄漏到生产包。
- 在示例和 smoke 测试中使用稳定 fixture，而不是连接真实生产系统。

不适用场景：业务模型库、生产 runtime client、隐式全局客户端封装、真实 adapter 默认实现、会自动读取生产密钥的库。

## 采用路径

1. 优先依赖已发布的 `github.com/ZoneCNH/testkitx` 版本；本地联调可临时使用 `replace` 指向 checkout，但发布前必须移除。
2. 只在测试文件、`testkit/`、测试工具、示例或显式 fixture 中 import `github.com/ZoneCNH/testkitx/pkg/testkitx/...`。
3. 在下游仓库加入生产 import 边界扫描，禁止 `pkg/*`、`internal/*` 等生产 Go 文件 import `testkitx`。
4. 使用 helper 产生的结构化 Evidence 记录 golden hash、contract hash、命令 stdout/stderr digest、manifest 字段和 fixture 路径。
5. 在下游仓库运行 `GOWORK=off go test ./...`、`GOWORK=off go vet ./...` 和该仓库自己的 CI gate。
6. 在采用声明中列出实际命令、CI artifact、依赖版本、commit SHA 和边界扫描结果。

示例：

```go
package example_test

import (
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/assertx"
	"github.com/ZoneCNH/testkitx/pkg/testkitx/fixture"
)

func TestConfigFixture(t *testing.T) {
	t.Parallel()

	workspace := fixture.NewWorkspace(t, "example.test/downstream")
	workspace.WriteOrFatal(t, "config.json", []byte(`{"name":"test"}`))

	assertx.Equal(t, "off", workspace.Env["GOWORK"])
}
```

## 采用检查清单

- [ ] 下游生产 Go 文件没有 import `github.com/ZoneCNH/testkitx/pkg/testkitx` 或其子包。
- [ ] 下游 module path 不包含 `x.go`，也没有 import `github.com/bytechainx/x.go`、`github.com/ZoneCNH/x.go` 或 `x.go/internal/*`。
- [ ] test helper 只服务测试、工具、示例或 fixture，不进入生产 runtime。
- [ ] golden 更新必须显式设置 `TESTKITX_UPDATE_GOLDEN=1`。
- [ ] command harness、contract、manifest 和 fixture Evidence 由下游仓库自己的 CI 或本地 gate 产生。
- [ ] 发布前移除临时 `replace`，并锁定实际依赖版本。

## 不足以作为采用证明的材料

- 只在 README 或计划文档中声明“已采用”。
- 只执行本地 dry-run，缺少命令输出或 CI artifact。
- 只生成 registry、manifest 草稿或 adoption plan。
- 依赖 `scripts/render_template.sh` 渲染出的生产模板，而没有证明下游只在测试路径使用 `testkitx`。

## 历史模板脚本

`scripts/render_template.sh` 仍保留为历史兼容和 integration regression 入口。它可以证明旧模板路径没有回归，但不是当前推荐的 L1 采用方式。新的下游采用应直接依赖 released `testkitx` helper，并用边界扫描证明它没有进入生产 import 图。
