"""配置优先级、密钥、migration。"""


def resolve(key: str, default: str, file: dict[str, str], env: dict[str, str]) -> tuple[str, str]:
    if key in env:
        return env[key], "env"
    if key in file:
        return file[key], "file"
    return default, "default"


def migrate(applied: set[int], all_versions: list[tuple[int, str]]) -> list[int]:
    ran = []
    for version, _ in all_versions:
        if version in applied:
            continue
        applied.add(version)
        ran.append(version)
    return ran


def main() -> None:
    file = {"db.name": "brm", "port": "9090"}
    env = {"port": "8080"}

    print("== 1. 配置优先级 ==")
    for k in ["port", "db.host", "db.name"]:
        value, origin = resolve(k, "localhost", file, env)
        print(f"{k} = {value}（来源：{origin}）")

    print("== 2. 密钥缺失拒绝启动 ==")
    secrets: dict[str, str] = {}
    if "DATABASE_URL" not in secrets:
        print("启动失败：DATABASE_URL 未设置（不要给密钥默认值）")

    print("== 3. migration ==")
    all_versions = [(1, "create accounts"), (2, "add email"), (3, "add index")]
    applied: set[int] = set()
    for v in migrate(applied, all_versions):
        print(f"应用 v{v} {all_versions[v - 1][1]}")
    print(f"当前版本 = {len(applied)}")
    again = migrate(applied, all_versions)
    print(f"再次执行：{len(again)} 个待应用，当前版本 = {len(applied)}")

    print("== 4. 结论 ==")
    print("配置分层要能说出「这个值从哪来」，否则排查时全靠猜")
    print("migration 必须可重复执行：靠版本表跳过，不靠人记得跑过没有")


if __name__ == "__main__":
    main()
