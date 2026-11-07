import { inspect } from 'node:util';

// —— class：构造器、公共字段、私有字段(#)、getter/setter、静态成员 ——
class User {
  #password;                       // 真正私有，外部访问会抛 SyntaxError
  static role = 'member';          // 静态字段属于类本身
  constructor(name, age, password) {
    this.name = name;
    this.age = age;
    this.#password = password;
  }
  get isAdult() { return this.age >= 18; }        // 像属性一样读
  set rename(v) { this.name = v; }                // 像属性一样写
  check(pwd) { return pwd === this.#password; }   // 私有字段只能在类内使用
  greet() { return `我是 ${this.name}`; }
}

const alice = new User('alice', 20, 's3cret');
console.log('实例方法：', alice.greet(), '，getter isAdult：', alice.isAdult);
console.log('静态字段：', User.role, '，私有校验：', alice.check('s3cret'), alice.check('wrong'));
console.log('外部看不到 #password，ownKeys =', inspect(Reflect.ownKeys(alice)));
alice.rename = 'bob';
console.log('setter 改名后：', alice.name);

// —— 继承 + super + instanceof ——
class Admin extends User {
  constructor(name, age, pwd, level) {
    super(name, age, pwd);
    this.level = level;
  }
  greet() { return `${super.greet()}，权限 ${this.level}`; }
}
const root = new Admin('root', 30, 'pw', 9);
console.log('继承后重写：', root.greet());
console.log('instanceof：', root instanceof Admin, root instanceof User, alice instanceof Admin);

// —— 没有接口：JS 靠鸭子类型（结构形状），TS 的 interface 运行时被擦除 ——
// TypeScript 里会写 interface Shape { area(): number }，编译成 JS 后什么都不剩
function totalArea(shapes) {
  return shapes.reduce((s, x) => s + x.area(), 0);
}
const shapes = [
  { area: () => 3 * 4 },
  { area: () => 5 * 5 },
  class Circle { constructor(r) { this.r = r; } area() { return Math.PI * this.r ** 2; } },
];
console.log('鸭子类型求和（前两项）：', totalArea(shapes.slice(0, 2)));
console.log('类实例同样满足形状：', totalArea([new shapes[2](2)]).toFixed(2));

// —— 「结构体」风格：用 Object.freeze 做不可变记录 ——
const Point = Object.freeze({ x: 1, y: 2, label: '原点' });
console.log('冻结记录：', JSON.stringify(Point), '，可扩展吗：', Object.isExtensible(Point));

// —— 模块：ESM 静态 import/export，顶层 this 是 undefined ——
console.log('ESM 顶层 this：', typeof this, '，模块 URL 末段：', import.meta.url.split('/').pop());
