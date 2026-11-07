import codecs


def main() -> None:
    print("1) 字符串不可变：索引赋值触发 TypeError")
    s = "abc"
    try:
        s[0] = "z"
    except TypeError:
        print("  s[0]='z' 抛 TypeError；要改只能新建字符串")
    print(f"  尝试后 s 仍是 {s!r}；'z'+s[1:]={'z' + s[1:]}")

    print("2) Unicode 码点 vs UTF-8 字节")
    text = "中文"
    print(f"  len('中文')={len(text)}（码点数）；len(encode)={len(text.encode('utf-8'))}（字节数）")
    print(f"  ord('中')={ord('中')}；chr(0x4E2D)={chr(0x4E2D)}；'中' 的 UTF-8 字节={list(text[0].encode('utf-8'))}")

    print("3) 遍历差异：按码点遍历 vs 按字节遍历")
    sample = "A中😀"
    print(f"  按码点 for: {[c for c in sample]}")
    print(f"  按字节前 6 个: {list(sample.encode('utf-8'))[:6]}")
    print(f"  码点数={len(sample)}；字节数={len(sample.encode('utf-8'))}")

    print("4) 拼接与格式化：join / f-string / format / %")
    parts = ["a", "b", "c"]
    name, age = "alice", 30
    legacy = "%s:%d" % (name, age)
    print(f"  join={'-'.join(parts)}；f-string={f'{name}:{age}'}；format={'{}:{}'.format(name, age)}；%= {legacy}")

    print("5) 转义与原始字符串")
    print(f"  转义制表符={repr('a\tb')}；换行符个数={'a\nb'.count(chr(10))}")
    print(f"  原始字符串={repr(r'a\nb')}（反斜杠不转义）")

    print("6) 大小写 / 查找 / 切分")
    word = "Hello,Python,World"
    print(f"  upper={word.upper()}；casefold={'Straße'.casefold()}；title={'hello world'.title()}")
    print(f"  find('Python')={word.find('Python')}；startswith={word.startswith('Hello')}；count(',')={word.count(',')}")
    print(f"  split={word.split(',')}；splitlines={'a\nb\n'.splitlines()}；partition={word.partition(',')}")

    print("7) 编码解码与错误策略")
    raw = "中文".encode("utf-8")
    print(f"  decode 还原={raw.decode('utf-8')}")
    broken = b"\xff\xfe"
    print(f"  errors='replace' 解码={broken.decode('utf-8', errors='replace')!r}")
    print(f"  codecs.BOM_UTF8 长度={len(codecs.BOM_UTF8)}")


if __name__ == "__main__":
    main()
