# ACCEPTANCE

本文件记录 `testkitx` 当前 release readiness 的验收标准。仓库只有在下面的事实同时成立时，才可对外宣称 ready for release。

## 功能验收

- `testkitx` 仍然是 L1 测试专用能力库，而不是生产运行时库。
- 公共 helper 只暴露在 `pkg/testkitx/`，并且保持测试、fixture、工具、示例的使用边界。
- contract / golden / harness / manifest 能力继续输出可审计 Evidence。
- `README.md`、`FEATURES.md` 和 `docs/release.md` 对功能边界和 release gate 的描述一致。

## 质量验收

- `GOWORK=off go test ./...` 通过。
- `GOWORK=off go vet ./...` 通过。
- `GOWORK=off make ci` 通过。
- `GOWORK=off make integration` 通过。
- `GOWORK=off make release-check` 通过。
- `GOWORK=off make release-final-check` 通过。
- `make release-preflight VERSION=vX.Y.Z` 在打 tag 前通过。

## 版本与文档验收

- `pkg/testkitx/version.go` 中的 module name 与仓库实际 module 保持一致。
- `CHANGELOG.md` 继续保留 `未发布` 入口，并记录当前已发布版本。
- `FEATURES.md`、`ACCEPTANCE.md` 和 `docs/release.md` 一致反映当前门禁、Evidence 和版本约束。

## 不接受项

- 生产包 import `github.com/ZoneCNH/testkitx/pkg/testkitx/...`
- 未显式设置 `TESTKITX_UPDATE_GOLDEN=1` 时更新 golden
- 在缺少 `golangci-lint` 或 `govulncheck` 时把必需 gate 当作可选
- 没有 Evidence artifact 仍声称 release ready

