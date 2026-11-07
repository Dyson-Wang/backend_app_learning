// —— 运行时类型信息 ——
class User {
  constructor(name) { this.name = name; }
  greet() { return `hi ${this.name}`; }
}
const u = new User('alice');
console.log('typeof / instanceof / constructor.name：', typeof u, u instanceof User, u.constructor.name);
console.log('Object.prototype.toString：', Object.prototype.toString.call([]), Object.prototype.toString.call(null));

// —— 按名字访问字段 / 方法（反射）——
const field = 'name';
console.log('u[field]：', u[field], '；Reflect.get：', Reflect.get(u, field));
Reflect.set(u, 'age', 20);
console.log('Reflect.set 后 age：', Reflect.get(u, 'age'), '；Reflect.has：', Reflect.has(u, 'age'));
const method = 'greet';
console.log('按名字调方法 Reflect.apply：', Reflect.apply(Reflect.get(u, method), u, []));
const created = Reflect.construct(User, ['bob']);
console.log('Reflect.construct 造对象：', created.greet(), '；是 User：', created instanceof User);

// —— 属性描述符 /「注解」的近似物 ——
// JS 没有原生注解/标签；Object.defineProperty 的属性描述符（writable/enumerable/configurable）
// 是最接近「属性/标签」的运行时元数据，可用 getOwnPropertyDescriptor 读出来。
const point = {};
Object.defineProperty(point, 'x', { value: 1, writable: false, enumerable: true, configurable: false });
Object.defineProperty(point, 'hidden', { value: 'secret', enumerable: false });
console.log('Object.keys（仅可枚举）：', Object.keys(point).join(','));
console.log('Reflect.ownKeys（含不可枚举）：', Reflect.ownKeys(point).join(','));
console.log('x 的描述符：', JSON.stringify(Object.getOwnPropertyDescriptor(point, 'x')));
let writeErr;
try { point.x = 9; } catch (err) { writeErr = err.constructor.name; }
console.log('不可写属性赋值后仍是：', point.x, '，严格模式报：', writeErr);

// —— Proxy：拦截基本操作，trap 内用 Reflect 转发默认行为 ——
const log = [];
const traced = new Proxy({ a: 1, b: 2 }, {
  get(target, key, receiver) { log.push(`get:${String(key)}`); return Reflect.get(target, key, receiver); },
  set(target, key, value, receiver) { log.push(`set:${String(key)}`); return Reflect.set(target, key, value, receiver); },
  has(target, key) { log.push(`has:${String(key)}`); return Reflect.has(target, key); },
});
traced.a; traced.a = 10; 'b' in traced;
console.log('Proxy get/set/has 拦截记录：', log.join(' '));
console.log('转发后真实值：', traced.a);

// 校验型 Proxy：set 之前做范围校验
function validated(obj) {
  return new Proxy(obj, {
    set(target, key, value, receiver) {
      if (key === 'age' && (!Number.isInteger(value) || value < 0)) {
        throw new RangeError('age 必须是非负整数');
      }
      return Reflect.set(target, key, value, receiver);
    },
  });
}
const profile = validated({ age: 18 });
profile.age = 30;
let validationErr;
try { profile.age = -1; } catch (err) { validationErr = `${err.constructor.name}:${err.message}`; }
console.log('校验 Proxy 合法值：', profile.age, '；非法值报：', validationErr);

// —— 边界：Proxy 不能违背对象不变式，也拦不到 # 私有字段 ——
const fixed = {};
Object.defineProperty(fixed, 'id', { value: 7, writable: false, configurable: false });
const liar = new Proxy(fixed, {
  get(target, key) { return key === 'id' ? 999 : Reflect.get(target, key); },
});
let liarErr;
try { liar.id; } catch (err) { liarErr = err.constructor.name; }
console.log('不可写不可配置属性不能被 Proxy 谎报：', liarErr, '（真实值', Reflect.get(fixed, 'id'), '）');
class Secret { #code = 42; getCode() { return this.#code; } }
const s = new Secret();
console.log('# 私有字段不在 ownKeys 里，长度：', Reflect.ownKeys(s).length, '；只能经方法读：', s.getCode());

// —— 装饰器 / 注解 / 宏：JS 运行时的现状 ——
// JS 无原生装饰器（TC39 Stage 3 推进中）、无注解、无宏、无编译期代码生成；
// TS 装饰器 + emitDecoratorMetadata 会在编译期产出 design:type 等元数据，
// 配合 reflect-metadata 运行时读取，但那是 tsc 产物，纯 .mjs 里跑不了。
console.log('原生 Reflect 方法齐全：', ['apply', 'construct', 'get', 'set', 'has', 'ownKeys'].every((k) => k in Reflect));
console.log('动态调用代价：Proxy/Reflect 破坏内联缓存、慢于直接调用，只用于框架/元编程层');
console.log('结论：JS 有运行时反射(Reflect/Proxy)，但没有注解/标签/宏，也没有编译期元编程');
