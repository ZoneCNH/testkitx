# Testkit 测试工具

为生成的基础库提供可复用测试夹具和断言。

## 契约

- `Config(name string)` 返回带 `Name` 和 `Timeout` 的最小有效配置。
- `RequireNoError(t, err)` 在 `err == nil` 时保持静默，在非空错误时终止当前测试。
- `RequireGolden(t, path, actual)` 读取 golden 文件并比较实际输出；不一致时报告 expected / actual 上下文。

## 回归覆盖

`fixture_test.go` 锁定 `Config("fixture")` 的字段和 `Validate` 结果，并验证 `RequireNoError(t, nil)` 可用。`golden_test.go` 锁定 golden 断言的匹配路径。生成后的基础库需要保留这组最小测试，以防测试夹具随包名替换、配置 contract 或稳定输出漂移。

生成的库应保持此包独立于 `x.go` 和业务特定模型。

## fixture.WriteOrFatal

`WriteOrFatal(t testing.TB, name string, data []byte)` 是推荐的 fixture 文件写入 API。它将 `data` 写入 workspace 中名为 `name` 的文件，写入失败时立即调用 `t.Fatal` 终止测试，确保测试不会在文件写入静默失败后继续执行。示例：

```go
workspace.WriteOrFatal(t, "config.json", []byte(`{"key":"value"}`))
```
