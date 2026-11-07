// 07 · 包管理与依赖
// 零第三方依赖：用标准库把 crate / mod / 路径 / 可见性 与 Cargo 元数据讲清楚。

// 一个 package 对应一个 Cargo.toml，默认产出一个同名 binary crate。
// mod 树是 crate 的内部命名空间；crate:: 指当前 crate 根。
mod inner {
    pub const LABEL: &str = "inner";

    pub fn value() -> i32 {
        10
    }

    pub mod deeper {
        // super:: 指父模块
        pub fn from_super() -> i32 {
            super::value() * 2
        }
    }

    // pub use 把深层项重新导出，调用方路径更短
    pub use self::deeper::from_super;
}

use inner::from_super;

fn main() {
    // 1) crate:: 绝对路径访问模块
    println!("crate::inner::LABEL = {}", crate::inner::LABEL);

    // 2) super:: 相对路径（在子模块里调用父模块函数）
    println!("super 相对路径 from_super = {}", from_super());

    // 3) pub use 重导出，调用方路径被缩短
    let v = from_super();
    println!("pub use 重导出结果 = {}", v);

    // 4) 平台相关常量由 Cargo 按目标注入
    println!("目标 os = {}, arch = {}", std::env::consts::OS, std::env::consts::ARCH);
    println!("cfg!(debug_assertions) = {}", cfg!(debug_assertions));

    // 5) env! 编译期读取 Cargo 注入的 CARGO_PKG_* 环境变量
    println!("包名(编译期) = {}, 版本 = {}", env!("CARGO_PKG_NAME"), env!("CARGO_PKG_VERSION"));

    // 6) std 本身就是一个外部 crate，通过 use 引入
    use std::collections::BTreeSet;
    let mut s = BTreeSet::new();
    s.insert("cargo");
    s.insert("crate");
    println!("std 外部 crate 示例: {:?}", s);
}
