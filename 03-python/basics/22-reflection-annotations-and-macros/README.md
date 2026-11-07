# 22 · Python · 反射、注解、宏与元编程

## 这个例子讲什么

运行时类型信息（`type` / `__annotations__` / `get_type_hints`）、按名字访问（`getattr`/`setattr`/`hasattr`）、`dataclasses.fields` 与 `inspect.signature`、装饰器 + `functools.wraps`、元类 `type`、`type()` 动态建类，以及动态调用的代价与边界。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/22-reflection-annotations-and-macros/main.py
```

## 关键点

- Python **一切皆对象**：类、函数、模块都能在运行时被检查；`__annotations__` 存原始注解，`get_type_hints()` 会解析成真实类型对象。
- `getattr(obj, "name")` 按字符串取属性，是 ORM / 序列化框架的基石；配合 `setattr` / `hasattr` / `callable` 做动态调用。
- `dataclasses.fields()` 拿到字段名、类型、默认值；`inspect.signature()` 读参数签名（装饰器要用 `@functools.wraps` 才不被 wrapper 遮住）。
- **装饰器**是「接收函数返回函数」的运行时替换；**元类**（`type` 子类）拦截类创建；`type("X", bases, ns)` 可动态造类——这些就是 Python 的「元编程」。
- **没有宏**：Python 不做编译期代码生成；反射走字典查找、比直接调用慢，且类型注解默认不强制（见 02 课）。

## 和其他语言的差异

- **Go**：`reflect` 包 + struct tag + `go:generate` 生成代码，**语言无宏**。
- **Node.js/TS**：`Reflect` / `Proxy`、TS 装饰器 + `emitDecoratorMetadata`（反射元数据）。
- **Java**：`java.lang.reflect` + 运行时注解 `RetentionPolicy.RUNTIME`，Lombok / APT 属编译期处理。
- **Rust**：无运行时反射，`std::any::Any` / `TypeId` + `macro_rules!` / 过程宏做编译期元编程。
- **Python**：`inspect` / `getattr` / `__annotations__` / `dataclasses.fields` / 装饰器 / 元类，无宏、注解不强制。
