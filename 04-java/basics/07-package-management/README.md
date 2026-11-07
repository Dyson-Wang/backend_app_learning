# 07 · Java · 包管理与依赖

## 这个例子讲什么

`classpath` 是什么、类加载器双亲委派、`Class.forName` 按名加载、依赖缺失抛 `ClassNotFoundException`、把 class 当资源读、包与 JPMS 模块边界。

## 怎么跑

```
cd backend-roadmap/04-java/basics/07-package-management
javac Main.java && java Main

# 带外部 jar 时：
# javac -cp lib/xxx.jar Main.java && java -cp .:lib/xxx.jar Main
# 本机无真实 JDK，示例未经编译验证
```

## 关键点

- **Java 语言本身没有包管理器**：依赖解析由构建工具（Maven 的 `pom.xml`、Gradle 的 `build.gradle`）完成，运行时靠 **classpath** 找 `.class` / `.jar`。
- **目录即包**（`package a.b;` 必须对应 `a/b/`），全限定类名 = 包 + 类名，是命名冲突的根本解法。
- 类加载器按 **双亲委派**：bootstrap(用 `null` 表示) ← platform ← app，先让父加载器尝试，避免用户类冒充 JDK 类。
- `Class.forName("...")` 是运行期按名字加载，找不到就是 `ClassNotFoundException`——这正是依赖没进 classpath 的典型报错。
- **JPMS**（Java 9）用 `module-info.java` 的 `requires/exports` 做强封装，运行在 `--module-path` 上；单文件示例处于**未命名模块 / 默认包**。

## 和其他语言的差异

| 能力 | Java | Go | Node.js | Python | Rust |
|---|---|---|---|---|---|
| 依赖清单 | `pom.xml`/`build.gradle` | `go.mod` | `package.json` | `requirements.txt`/`pyproject.toml` | `Cargo.toml` |
| 包管理命令 | mvn / gradle | `go mod` | npm / pnpm | pip / poetry | cargo |
| 依赖查找 | classpath / module-path | module cache | `node_modules` 逐级向上 | `sys.path` / site-packages | crate 图 + 锁文件 |
| 命名空间 | `package` + 全限定名 | import path（URL） | 模块 + 包名 | import 路径 + 虚拟环境 | `crate::mod` |
| 模块强封装 | JPMS `exports` | 无 | ESM `export` | 约定（无强制） | `pub` / `mod` |
