# xlib-standard 同步说明

## 结论

`testkitx` 需要同步 `xlib-standard` 中适用于 L1 测试库的共享标准，但不应做全量同步。同步目标是让本仓库继承基础库体系的安全、toolchain、Evidence 和边界约束，同时保持 `testkitx` 的 test-only 身份。

本次参考基线为 `xlib-standard` 的 `origin/main`，对应 tag `v0.4.13`。若本地 `xlib-standard/main` 与远端基线分叉，应以远端基线为准，避免把未确认的本地变更同步进来。

## 应同步范围

- Docker toolchain contract：`Dockerfile`、`.dockerignore`、`docker-compose.yml`、`.devcontainer/devcontainer.json`、`scripts/docker/*` 和 Makefile 中的 Docker gate 入口。
- 安全 gate：Security workflow 的 pull request、manual dispatch、weekly schedule、tool pinning、Go module version source 和 `GOWORK=off` 运行约束。
- Contract schema：Docker toolchain schema 和 downstream adoption proof schema。
- 文档口径：身份、采用、历史生成、release Evidence、下游采用证明和生产 import 边界说明。
- L1 兼容入口：必要的 Makefile alias 可以保留，但不得暗示存在生产二进制或 `goalcli` runtime。

## 不应同步范围

- 把 `testkitx` 恢复为生产基础库模板或默认下游生成源的叙述。
- 真实 runtime adapter、业务模型、生产连接、隐式密钥读取、全局 client 或 `x.go` 依赖。
- 与 L1 测试 helper 无关的应用层命令、服务部署、业务 contract 或下游专有 profile。
- 未经验证的本地 `xlib-standard` 分叉变更。

## 本仓库落地方式

`testkitx` 的 Docker 和 security 入口服务于本仓库自身的 test helper、contracts、历史模板 regression 和 release Evidence。Docker toolchain 只提供可重复的 Go 工具链容器，不引入生产运行时。

历史模板生成脚本继续作为 integration regression 使用；新的下游采用应直接依赖 released `testkitx` helper，并通过 `boundarytest` 或等价扫描证明生产 import 图没有依赖 `testkitx`。

## 验证口径

完成同步后至少验证：

- `GOWORK=off go test ./...`
- `GOWORK=off go vet ./...`
- `GOWORK=off make contracts`
- `GOWORK=off make drift-check`
- `bash -n scripts/docker/check_toolchain.sh scripts/docker/docker_gate.sh scripts/docker/check_contract.sh`

若本机没有 Docker daemon，可用 `XLIB_DOCKER_ALLOW_MISSING=1 make docker-toolchain-check` 验证缺失 Docker 时的报告路径；真实 Docker build 和 container gate 仍需在具备 Docker 的环境执行。
