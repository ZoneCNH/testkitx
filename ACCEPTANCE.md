# ACCEPTANCE

以下标准用于确认当前发布版本已经与 release gate、changelog 和版本元数据同步。

## Release acceptance

- `pkg/testkitx/version.go` 中的 `Version` 与 `.repo-contract.yaml` 的 `latest_git_tag` 一致。
- `CHANGELOG.md` 包含当前版本的标题行。
- `FEATURES.md` 和本文件都引用同一当前版本，不描述过期能力。
- `GOWORK=off go test ./...` 通过。
- `GOWORK=off go vet ./...` 通过。
- `GOWORK=off make release-check` 通过。
- `GOWORK=off make release-final-check` 通过。
- `make release-preflight VERSION=<current-version>` 在满足前置条件时通过。

## Current version

- 当前版本：`v0.4.1`
- 发布日期：`2026-06-21`
