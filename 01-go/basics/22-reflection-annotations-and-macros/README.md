# 22 · 反射、注解、宏与元编程

## 这个例子讲什么

用 `reflect` 包读运行时类型信息、遍历字段、按名字读写字段、按名字调用方法；用 struct tag 充当「注解」；说明 `go:generate` 是编译期代码生成的入口，以及 Go **没有宏**这一事实与动态调用的边界。

## 怎么跑

```bash
cd backend-roadmap/01-go
go run ./basics/22-reflection-annotations-and-macros
# 演示 //go:generate 指令（只打印，不真正生成）
go generate ./basics/22-reflection-annotations-and-macros
```

## 关键点

- **运行时类型信息**：`reflect.TypeOf` / `ValueOf`，`Kind()` 看底层类别（`struct`/`int`/`slice`…），`NumField`/`Field(i)` 遍历字段。
- **按名字访问**：`FieldByName("Name")`；要修改必须传指针并 `Elem()`，且 `CanSet()` 为真。`ValueOf(u)`（值）不可设，`ValueOf(&u).Elem()`（可寻址）才可设。
- **struct tag 就是 Go 的「注解」**：反引号里的 `json:"..." validate:"..."` 只是字符串，靠 `StructTag.Get` 按 key 读取，编译器不解释、运行时由库（如 `encoding/json`）解析。没有真正的自定义注解类型、没有保留策略。
- **动态调用**：`MethodByName("Greet").Call([]reflect.Value{...})`，参数与返回值都是 `[]reflect.Value`，类型错了**运行时 panic**，没有编译期检查。
- **边界**：未导出字段 `CanSet()==false`（本例 `age`），强设会 panic；反射慢、丢失类型安全、无法内联，只适合框架/序列化/ORM 等通用代码。
- **元编程**：Go **没有宏、没有装饰器、没有编译期 AST 变换**。`//go:generate` 指令配合 `stringer`/`mockgen` 等外部命令在编译前生成 `.go` 源文件，是「代码生成」而非宏；本例的 `//go:generate` 只是占位演示。

## 和其他语言的差异

| 语言 | 运行时反射 | 注解 / 属性 | 编译期元编程 |
|---|---|---|---|
| Go | `reflect`（Type/Value/Kind） | struct tag（纯字符串，库自行解析） | `go:generate` 代码生成；**无宏** |
| Node.js | `Reflect` / `Proxy` 可拦截操作 | 装饰器（TS），JS 里靠 `Reflect.metadata` | Babel/TS 转换器；JS 无原生宏 |
| Python | `inspect`、`__annotations__`、`dataclasses.fields` | 装饰器（运行时生效）、类型注解 | 装饰器/元类 `__metaclass__`；无宏 |
| Java | `java.lang.reflect`，可访问私有（setAccessible） | 注解 + `RetentionPolicy.RUNTIME`/`SOURCE` | Lombok / APT 注解处理器生成字节码或源码 |
| Rust | 无运行时反射，`std::any::Any` + `TypeId` 仅够类型识别 | `#[derive]` / 属性宏（过程宏） | `macro_rules!` + 过程宏，编译期展开、零运行时开销 |

把这一课的结论补到 `../../../99-cross-language/scenario-matrix.md` 的基础知识对照区。
