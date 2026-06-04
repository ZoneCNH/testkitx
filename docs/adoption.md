# 采用说明

## 适用场景

当需要创建新的 x.go 基础设施基础库时，优先采用 `testkitx` 作为模板源。适用对象包括但不限于：

- `foundationx`
- `configx`
- `observex`
- `postgresx`
- `kafkax`
- `redisx`
- `taosx`
- `ossx`

不适用场景：业务模型库、交易语义库、隐式全局客户端封装、会自动读取生产密钥的库。

## 采用路径

1. 从干净的 `testkitx` 工作区开始，确认基础 gate 可运行。
2. 使用 `scripts/render_template.sh` 渲染目标库，而不是手工批量替换。
3. 在生成后的目标库内运行 `GOWORK=off go test ./...`。
4. 根据目标库 profile 补充具体 adapter、配置字段、health check 和 integration smoke；profile 可参考 `docs/test-strategy.md` 中的 Pure、Config、Observability、Storage 和 Messaging 分层。
5. 运行目标库自己的 `make ci`、`make integration` 和 `make release-check`。
6. 检查目标库生成的 `release/manifest/latest.json`，确认 module、commit、tree SHA、contract 指纹、工具版本和 gate 结果来自目标库自身。
7. 在发布或合并声明中使用 `DONE with evidence:`，列出实际命令和 artifact。

示例：

```bash
scripts/render_template.sh \
  --module-name foundationx \
  --module-path github.com/ZoneCNH/foundationx \
  --package-name foundationx \
  --out ../foundationx

cd ../foundationx
GOWORK=off go test ./...
GOWORK=off make release-check
```

## 采用检查清单

- [ ] 目标 module path 不包含 `x.go`。
- [ ] 目标库未 import `github.com/bytechainx/x.go` 或 `github.com/ZoneCNH/x.go`。
- [ ] 目标库没有业务模型、生产连接默认值或真实密钥。
- [ ] `Config.Validate` 和 `Config.Sanitize` 已按目标库字段更新。
- [ ] `HealthCheck`、错误分类和 metrics contract 与目标库语义一致。
- [ ] contracts 与公共常量、JSON 字段和文档一致。
- [ ] examples 可以作为 smoke 测试运行。
- [ ] release Evidence 由目标库生成，不复用 `testkitx` 的 `latest.json`。

## 采用风险

- 手工复制模板容易遗漏包名、module path、imports、contract 指纹和 manifest 字段，必须优先使用生成脚本。
- 直接修改历史 Goal 文档可能掩盖当前实现事实；采用者应以当前代码、contracts 和 gate 输出为准。
- 如果目标库需要真实基础设施 integration，密钥路径只能由调用方显式传入，不能写入模板默认值。
