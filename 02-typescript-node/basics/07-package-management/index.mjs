import { builtinModules } from 'node:module';
import process from 'node:process';

// —— node: 前缀：强制解析到内置模块，杜绝与同名 npm 包混淆 ——
console.log('node:path 解析为：', import.meta.resolve('node:path'));
console.log('node:fs   解析为：', import.meta.resolve('node:fs'));
console.log('无前缀 fs 也会被规范化成：', import.meta.resolve('fs'));

// 未安装的第三方包解析直接抛错——本仓库零依赖，不 import 任何外来包
try {
  import.meta.resolve('left-pad');
} catch (err) {
  console.log('解析未安装的包抛错：', err.code);
}

// —— builtinModules：内置模块清单（名字不带 node: 前缀）——
console.log('内置模块总数：', builtinModules.length);
console.log('清单含 fs / worker_threads：', builtinModules.includes('fs'), builtinModules.includes('worker_threads'));
console.log('清单里带 node: 前缀的条目数：', builtinModules.filter((m) => m.startsWith('node:')).length);

// —— process.versions：编译进二进制的运行时与各依赖版本表 ——
console.log('node：', process.versions.node, '，v8：', process.versions.v8, '，libuv：', process.versions.uv);

// —— 包管理概念（讲「怎么组织」，不联网、不落盘、不创建 package.json）——
// package.json："type":"module" 让 .js 也按 ESM 解析；"exports" 限制对外暴露的子路径
// dependencies / devDependencies / peerDependencies / optionalDependencies 四类依赖
// semver：^1.2.3 放开 minor+patch，~1.2.3 只放开 patch，精确版本最可复现
// package-lock.json 锁定整棵依赖树：npm ci 严格按锁装，npm install 可能更新锁
// npm ls 看依赖树、npm audit 查漏洞；node_modules 的查找是逐级向上的目录回溯
const manifest = { name: 'demo', type: 'module', exports: { '.': './index.mjs' }, dependencies: {} };
console.log('示意 manifest（仅打印不写文件）：', JSON.stringify(manifest));
console.log('本课只用内置模块：无需 npm install，无需联网，无 package.json');
