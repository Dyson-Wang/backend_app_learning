// JS 没有原生 enum / 模式匹配。用 Object.freeze 当枚举常量的容器，
// 用「判别联合」（带 tag 字段的对象）+ switch 模拟 ADT 与穷尽匹配。
// TypeScript 里对应的写法是：
//   const Direction = { Up: 'UP', Down: 'DOWN' } as const;
//   type Direction = typeof Direction[keyof typeof Direction];
//   type Shape =
//     | { kind: 'circle'; r: number }
//     | { kind: 'rect'; w: number; h: number }
//     | { kind: 'text'; value: string };
// 但 .mjs 直接跑的是 JS，类型注解会被擦除，只能靠运行时 tag 自己判断。

// —— 枚举定义与取值：Object.freeze 防止误改 ——
const Direction = Object.freeze({
  Up: 'UP', Down: 'DOWN', Left: 'LEFT', Right: 'RIGHT',
});
console.log('枚举取值：', Direction.Up, Direction.Right);
console.log('枚举全部值：', Object.values(Direction).join(','));
const frozen = Object.isFrozen(Direction);
let afterWrite;
try { Direction.Up = 'BAD'; } catch { /* 严格模式对冻结对象赋值抛 TypeError */ }
afterWrite = Direction.Up;
console.log('冻结了吗：', frozen, '；尝试改写后仍为：', afterWrite);

// 值到名字的反向查找要自己建表（TS 的 enum 会自动生成反向映射）
const reverse = Object.fromEntries(Object.entries(Direction).map(([k, v]) => [v, k]));
console.log('反向查找 UP =>', reverse[Direction.Up]);

// —— 带数据的枚举 / ADT：每个分支是带 tag 的普通对象 ——
const shapes = [
  { kind: 'circle', r: 2 },
  { kind: 'rect', w: 3, h: 4 },
  { kind: 'text', value: 'hi' },
];

// 穷尽性：switch 覆盖所有 tag，default 用运行时抛错兜底。
// TS 里会写 const _exhaustive: never = s; 让编译器在漏分支时报错。
function area(s) {
  switch (s.kind) {
    case 'circle': {
      const { r } = s;                 // 解构：从分支里取数据
      return Math.PI * r * r;
    }
    case 'rect': {
      const { w, h } = s;
      return w * h;
    }
    case 'text': {
      const { value } = s;
      return value.length;
    }
    default:
      throw new Error(`未处理的分支：${s.kind}`); // 运行时穷尽兜底
  }
}
console.log('ADT 面积/长度：', shapes.map((s) => area(s).toFixed(2)).join(' / '));

// 漏分支在 JS 里不会编译报错，只会被 default 兜住——这正是「局限」
let caught = '';
try { area({ kind: 'triangle', a: 1 }); } catch (err) { caught = err.message; }
console.log('未知分支被 default 接住：', caught);

// —— switch(true) 做「guard」式分支（模拟带条件的模式）——
function describe(s) {
  switch (true) {
    case s.kind === 'circle' && s.r > 5: return '大圆';
    case s.kind === 'circle': return '小圆';
    case s.kind === 'rect': return '矩形';
    default: return '其它';
  }
}
console.log('switch(true) 守卫式匹配：', describe(shapes[0]), describe(shapes[1]));

// TS 穷尽性检查写法（仅注释，.mjs 不生效）：
//   function assertNever(x: never): never { throw new Error(String(x)); }
// 在 default 分支调用 assertNever(s)，漏了任一 kind 就编译不过。
console.log('局限：JS 无解构模式/守卫/类型收窄，穷尽性只能靠运行时 default 或 TS 编译期 never');
