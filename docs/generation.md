# 历史模板生成

## 当前定位

`scripts/render_template.sh` 是历史模板兼容和 integration regression 入口。它用于证明旧的模板渲染路径、包名替换、contract gate、boundary gate 和 release Evidence 生成没有意外回归。

当前 `testkitx` 的主身份是 L1 测试专用能力库。新的下游基础库采用 `testkitx` 时，应直接在测试代码、测试工具、示例或 fixture 中 import released helper，并通过边界扫描证明生产 import 图没有依赖 `testkitx`。不要把模板渲染作为默认采用路径。

## 兼容用途

仍可在以下场景使用生成脚本：

- 维护历史模板资产的 regression 测试。
- 验证包名、module path、imports、contract 和 manifest 替换逻辑。
- 临时构造下游 fixture，用于检查当前仓库 gate 是否覆盖旧模板路径。

不应使用生成脚本来证明新的 L1 test-only 采用完成。采用证明必须来自真实下游仓库或明确的测试 fixture，并包含生产 import 边界扫描结果。

## 示例

```bash
scripts/render_template.sh \
  --module-name foundationx \
  --module-path github.com/ZoneCNH/foundationx \
  --package-name foundationx \
  --out ../foundationx
```

`--out` 必须指向不存在或为空的目录，避免覆盖已有仓库内容。

## 渲染范围

- `testkitx` 替换为 `--module-name`。
- `github.com/ZoneCNH/testkitx` 替换为 `--module-path`。
- `pkg/testkitx`、包名和 imports 替换为 `--package-name`。
- 文档、Go 代码、JSON contract、shell 脚本、Makefile 和 CI 配置同步更新。

脚本不会复制 `.git`、`.omc`、`.omx`、`.worktree`、`release/manifest/latest.json` 和 `release/manifest/latest.json.sha256`。这些 release Evidence 文件是生成产物，生成后的库必须自己运行 release gate 生成新的 Evidence artifact。

## 回归验证

模板自身的 `make integration` 会渲染两个临时下游库：

- `foundationx`：目标仓库路径 `github.com/ZoneCNH/foundationx`，用于证明真实迁移目标仍可生成。
- `corekit`：中性路径 `example.com/acme/corekit`，用于证明替换逻辑不依赖特定组织或包名。

每个临时库都会运行以下验证：

- `scripts/check_rendered_template.sh`
- `GOWORK=off go test ./...`
- `GOWORK=off make contracts`
- `GOWORK=off make boundary`
- `CHECK_STATUS=passed GOWORK=off make evidence`
- `RELEASE_EVIDENCE_REQUIRE_PASSED=1 GOWORK=off make release-evidence-check`

这组验证仅代表历史模板 regression 通过，不替代 `testkitx` 作为 L1 测试库的下游采用证明。
