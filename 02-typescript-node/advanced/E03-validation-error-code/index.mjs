// 统一错误结构：code 给程序判断，message 给人看
const Code = Object.freeze({
  REQUIRED: 'REQUIRED',
  LENGTH: 'LENGTH',
  FORMAT: 'FORMAT',
  RANGE: 'RANGE',
});

function validate(user) {
  const errors = [];

  const name = (user.name ?? '').trim();
  if (name === '') {
    errors.push({ code: Code.REQUIRED, field: 'name', message: 'name 不能为空' });
  } else if (name.length < 2 || name.length > 20) {
    errors.push({ code: Code.LENGTH, field: 'name', message: 'name 长度必须在 2-20 之间' });
  }

  const email = user.email ?? '';
  if (email === '') {
    errors.push({ code: Code.REQUIRED, field: 'email', message: 'email 不能为空' });
  } else if (!email.includes('@')) {
    errors.push({ code: Code.FORMAT, field: 'email', message: 'email 格式不正确' });
  }

  const age = user.age;
  if (age === undefined || age === null) {
    errors.push({ code: Code.REQUIRED, field: 'age', message: 'age 不能为空' });
  } else if (!Number.isInteger(age) || age < 0 || age > 150) {
    errors.push({ code: Code.RANGE, field: 'age', message: 'age 必须在 0-150 之间' });
  }

  return errors;
}

function report(label, user) {
  const errors = validate(user);
  console.log(`${label}：${errors.length} 个错误`);
  for (const e of errors) {
    console.log(`  [${e.code}] ${e.field}: ${e.message}`);
  }
}

report('全空', {});
report('短名 + 坏邮箱 + 超龄', { name: 'a', email: 'bad', age: 200 });
report('合法', { name: 'ada', email: 'ada@example.com', age: 36 });

console.log('结论：一次返回全部错误，前端可逐字段标红；code 与 message 分离便于国际化');
