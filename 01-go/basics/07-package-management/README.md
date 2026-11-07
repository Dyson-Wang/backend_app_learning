# 07 · 包管理与依赖

## 这个例子讲什么

用真实可编译的 `internal/` 子包演示 import 路径与可见性，并讲清 `go.mod` 的每条指令。本仓库零第三方依赖。

## 怎么跑

```
cd backend-roadmap/01-go

go run ./basics/07-package-management
go list -m all                    # 列出模块依赖（这里只有自己）
go list ./...                     # 列出本模块所有包
go mod tidy                       # 增删依赖，清理 go.sum
go build ./...                    # 一次编译全部包
```

## 关键点

- **import 路径 = module 路径 + 目录相对路径**。本模块是 `github.com/dyson/backend-roadmap/go`，子包就是 `.../go/internal/toolkit`。
- `internal/` 是编译器强制的可见性边界：只有 `internal` 父目录内的包能 import 它。适合放不愿暴露给外部的实现细节。
- `go.mod` 三件套：`module`（模块路径）、`go`（最低语言版本）、`require`（依赖及版本）。`go.sum` 记录依赖哈希，**必须提交**。
- `replace` 用于本地联调或替换有问题的依赖：`replace example.com/lib => ../local-lib`；`exclude` 排除特定版本。
- `go mod tidy` 会移除没被引用的依赖；`go mod vendor` 把依赖拷进 `vendor/` 以便离线构建。
- Go 的**最小版本选择（MVS）**：取所有依赖要求中的最大版本，不做回溯，因此依赖解析可预测。

`advanced/` 里引入第三方库时（如 `pgx`、`redis`、`prometheus`），会在对应目录单独给依赖清单，不污染这里的零依赖基线。

## 和其他语言的差异

| 语言 | 清单文件 | 锁文件 | 隔离方式 |
|---|---|---|---|
| Go | `go.mod` | `go.sum` | 模块缓存，MVS |
| Node.js | `package.json` | `pnpm-lock.yaml` | `node_modules` |
| Python | `requirements.txt` / `pyproject.toml` | 无统一标准 | `venv` |
| Java | `pom.xml` / `build.gradle` | 依赖解析缓存 | Maven/Gradle 本地仓库 |
| Rust | `Cargo.toml` | `Cargo.lock` | cargo registry + target/ |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
