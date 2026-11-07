import { readFileSync, writeFileSync, createReadStream, mkdtempSync, rmSync } from 'node:fs';
import { readFile, writeFile } from 'node:fs/promises';
import { join } from 'node:path';
import { tmpdir } from 'node:os';
import { fileURLToPath } from 'node:url';
import { TextDecoder } from 'node:util';

// 在系统临时目录里造文件，跑完即删，不污染仓库
const dir = mkdtempSync(join(tmpdir(), 'brm-08-'));
const file = join(dir, 'data.json');
try {
  // —— 同步 API：脚本方便，但会阻塞事件循环，服务里禁用 ——
  writeFileSync(file, '{"name":"alice"}', 'utf8');
  console.log('同步写入后读回：', readFileSync(file, 'utf8'));

  // —— 异步 promises API：生产推荐 ——
  await writeFile(file, JSON.stringify({ name: 'bob', tags: ['dev'] }), 'utf8');
  const text = await readFile(file, 'utf8');
  console.log('异步读回：', text);

  // —— JSON：parse 的 reviver / stringify 的 replacer 与缩进 ——
  const data = JSON.parse(text, (k, v) => (k === 'name' ? String(v).toUpperCase() : v));
  console.log('reviver 改写 name：', data.name);
  const pretty = JSON.stringify(data, (k, v) => (k === 'tags' ? v.join('/') : v), 2);
  console.log('replacer + 缩进：\n' + pretty);
  console.log('缩进后行数：', pretty.split('\n').length);

  // JSON 的能力边界：BigInt / undefined / 循环引用都不能直接序列化
  try {
    JSON.stringify({ big: 1n });
  } catch (err) {
    console.log('BigInt 序列化报错：', err.constructor.name);
  }
  console.log('undefined 字段被忽略：', JSON.stringify({ a: undefined, b: 2 }));

  // —— Buffer 与编码：字节层 vs 字符层 ——
  const buf = Buffer.from('中文', 'utf8');
  console.log('「中文」UTF-8 字节数：', buf.length, '，十六进制：', buf.toString('hex'));
  console.log('TextDecoder 解码回字符：', new TextDecoder('utf8').decode(buf));
  console.log('字节长度 vs 字符长度：', Buffer.byteLength('中文'), '中文'.length);

  // —— 流式读取：for await 消费 createReadStream，适合大文件 ——
  let streamed = '';
  for await (const chunk of createReadStream(file)) streamed += chunk;
  console.log('流式拼接字节数：', Buffer.byteLength(streamed), '，与 readFile 一致：', streamed === text);

  // —— 模块文件定位：import.meta.url -> 真实文件路径 ——
  console.log('本入口文件名：', fileURLToPath(import.meta.url).split('/').pop());
} finally {
  rmSync(dir, { recursive: true, force: true }); // 清理临时目录
  console.log('临时目录已清理完成');
}
