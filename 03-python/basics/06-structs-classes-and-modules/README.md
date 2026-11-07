# 06 · Python · 结构体 / 类、接口、模块

## 这个例子讲什么

`@dataclass` 自动生成特殊方法、`@property` 计算属性、`abc.ABC` 抽象基类作为名义接口、`typing.Protocol` 结构化接口、类属性 vs 实例属性、`__len__` 鸭子类型，以及「每个 `.py` 文件即模块」。

## 怎么跑

```
cd backend-roadmap/03-python
python3 basics/06-structs-classes-and-modules/main.py
```

## 关键点

- **`@dataclass`**：按注解自动生成 `__init__/__repr__/__eq__`，`field(default_factory=list)` 避开可变默认值坑（见 04 课）。
- **两种「接口」**：`abc.ABC` + `@abstractmethod` 是**名义**接口（必须继承并实现）；`Protocol` 是**结构化**接口（鸭子类型，无需继承，14 课详述）。
- `abc.ABC` 带未实现抽象方法时**不能实例化**，调用即 `TypeError`。
- 类属性写在类体、被所有实例共享；实例属性在 `__init__` 里各自独立。
- Python 没有 `public/private` 关键字，可见性靠约定（17 课）：`_x` 私有约定、`__x` 名称改写。

## 和其他语言的差异

- **Go**：`struct` + 方法接收者、隐式满足的 `interface`（结构化）、`package` 目录即包。
- **Node.js/TS**：`class` / `interface`（TS 编译期擦除）、ESM `export/import`、`#private` 真私有字段。
- **Java**：`class` / `interface` / `record`、`extends` / `implements`、`public/private` 关键字、包即目录。
- **Rust**：`struct` / `enum` + `impl`、`trait`（可 `dyn` 或泛型）、`mod` / `crate` 模块树。
- **Python**：`class` + `@dataclass`、`abc.ABC` 抽象方法、`Protocol` 结构化接口、单文件即模块、无访问控制关键字。
