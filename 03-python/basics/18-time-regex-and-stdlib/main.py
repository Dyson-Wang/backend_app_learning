import functools
import math
import random
import re
import statistics
from collections import Counter
from datetime import datetime, timedelta, timezone


def main() -> None:
    print("1) 时间点 / 时长 / 时区：用固定时间，不打印当前绝对时间")
    cst = timezone(timedelta(hours=8))
    utc_noon = datetime(2024, 1, 1, 12, 0, tzinfo=timezone.utc)
    local = utc_noon.astimezone(cst)
    print(f"  UTC 12:00 转 +08:00 = {local:%Y-%m-%d %H:%M %z}")
    delta = datetime(2024, 1, 3, 15, 0) - datetime(2024, 1, 1, 12, 0)
    print(f"  两个固定时刻差值={delta}（天数={delta.days}）")
    print(f"  fromisoformat={datetime.fromisoformat('2024-01-01T08:00:00').isoformat()}")

    print("2) 正则：匹配 / 捕获 / 替换 / 切分")
    pat = re.compile(r"(?P<user>\w+)@(?P<host>[\w.]+)")
    m = pat.search("联系 alice@example.com 或 bob@test.org")
    print(f"  命名组: user={m.group('user') if m else None} host={m.group('host') if m else None}")
    print(f"  findall={re.findall(r'\d+', 'a1b22c333')}；sub={re.sub(r'\d+', '#', 'a1b22')}")
    print(f"  切分={re.split(r'[,;]\s*', 'a, b; c')}；忽略大小写匹配={bool(re.match('abc', 'ABC', re.IGNORECASE))}")

    print("3) 随机数：固定种子保证可复现")
    rnd = random.Random(42)
    rolls = [rnd.randint(1, 6) for _ in range(5)]
    print(f"  Random(42) 五次掷骰={rolls}；choice={rnd.choice(['a','b','c'])}")
    rnd2 = random.Random(42)
    print(f"  同种子可复现={[rnd2.randint(1, 6) for _ in range(5)] == rolls}")

    print("4) 排序与常用工具：sorted / Counter / statistics / reduce / math")
    words = ["banana", "fig", "apple", "kiwi"]
    print(f"  sorted(key=len)={sorted(words, key=len)}；Counter={Counter('abracadabra').most_common(2)}")
    print(f"  statistics.mean={statistics.mean([1,2,3,4])}；reduce 阶乘={functools.reduce(lambda a,b:a*b, range(1,5))}")
    print(f"  math.isqrt(50)={math.isqrt(50)}；pow(2,10,1000)={pow(2,10,1000)}")

    print("5) timedelta 运算与解析")
    shifted = local + timedelta(days=10, hours=-3)
    print(f"  local + 10天-3小时 = {shifted:%Y-%m-%d %H:%M %z}")
    print(f"  strptime('2024-02-29')={datetime.strptime('2024-02-29','%Y-%m-%d').date()}")


if __name__ == "__main__":
    main()
