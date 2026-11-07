import json
import tempfile
from pathlib import Path


def main() -> None:
    # 用 TemporaryDirectory，退出即清理，不往仓库里留文件
    with tempfile.TemporaryDirectory(prefix="brm08_") as tmp:
        root = Path(tmp)
        print("1) pathlib 写文本 / 读文本 / with 自动关闭")
        note = root / "note.txt"
        note.write_text("第一行：中文\n第二行：hello\n", encoding="utf-8")
        print(f"  写入字节数={note.stat().st_size}；read_text 行数={len(note.read_text(encoding='utf-8').splitlines())}")

        print("2) 二进制写与 UTF-8 编码：字符数 != 字节数")
        raw = root / "raw.bin"
        text = "中文"
        payload = text.encode("utf-8")
        raw.write_bytes(payload)
        back = raw.read_bytes().decode("utf-8")
        print(f"  '中文' 字符数={len(text)}；UTF-8 字节数={len(payload)}；解码回={back}")

        print("3) json 序列化：dumps / loads / ensure_ascii")
        data = {"name": "爱丽丝", "age": 30, "tags": ["a", "b"]}
        compact = json.dumps(data, ensure_ascii=False, sort_keys=True)
        print(f"  dumps(ensure_ascii=False)={compact}")
        print(f"  loads 还原 name={json.loads(compact)['name']}")

        print("4) 直接读写文件对象：json.dump / json.load")
        fp = root / "data.json"
        with fp.open("w", encoding="utf-8") as f:
            json.dump(data, f, ensure_ascii=False)
        with fp.open("r", encoding="utf-8") as f:
            loaded = json.load(f)
        print(f"  json.load 结果相等={loaded == data}；文件存在={fp.is_file()}")

        print("5) 追加写与 open 的 'a' 模式")
        with (root / "log.txt").open("w", encoding="utf-8") as f:
            f.write("start\n")
        with (root / "log.txt").open("a", encoding="utf-8") as f:
            f.write("end\n")
        print(f"  追加后内容={(root / 'log.txt').read_text(encoding='utf-8').strip()!r}")

        print("6) 错误处理：非法 JSON 抛 JSONDecodeError")
        try:
            json.loads("{not json")
        except json.JSONDecodeError as exc:
            print(f"  捕获 JSONDecodeError: msg 非空={bool(exc.msg)}；pos={exc.pos}")
        print(f"  当前临时目录存在={root.exists()}（with 退出后自动清理）")


if __name__ == "__main__":
    main()
