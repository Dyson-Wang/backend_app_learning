from enum import Enum


class Code(Enum):
    REQUIRED = "REQUIRED"
    LENGTH = "LENGTH"
    FORMAT = "FORMAT"
    RANGE = "RANGE"


def validate(user):
    errors = []

    name = (user.get("name") or "").strip()
    if not name:
        errors.append({"code": Code.REQUIRED.value, "field": "name", "message": "name 不能为空"})
    elif not (2 <= len(name) <= 20):
        errors.append({"code": Code.LENGTH.value, "field": "name", "message": "name 长度必须在 2-20 之间"})

    email = user.get("email") or ""
    if not email:
        errors.append({"code": Code.REQUIRED.value, "field": "email", "message": "email 不能为空"})
    elif "@" not in email:
        errors.append({"code": Code.FORMAT.value, "field": "email", "message": "email 格式不正确"})

    age = user.get("age")
    if age is None:
        errors.append({"code": Code.REQUIRED.value, "field": "age", "message": "age 不能为空"})
    elif not isinstance(age, int) or isinstance(age, bool) or not (0 <= age <= 150):
        errors.append({"code": Code.RANGE.value, "field": "age", "message": "age 必须在 0-150 之间"})

    return errors


def report(label, user):
    errors = validate(user)
    print(f"{label}：{len(errors)} 个错误")
    for e in errors:
        print(f"  [{e['code']}] {e['field']}: {e['message']}")


def main():
    report("全空", {})
    report("短名 + 坏邮箱 + 超龄", {"name": "a", "email": "bad", "age": 200})
    report("合法", {"name": "ada", "email": "ada@example.com", "age": 36})
    print("结论：一次返回全部错误，前端可逐字段标红；code 与 message 分离便于国际化")


if __name__ == "__main__":
    main()
