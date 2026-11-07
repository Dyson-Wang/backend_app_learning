"""JWT 刷新与吊销。"""
import hashlib
import hmac

SECRET = b"s3cr3t"


def sign(payload: str) -> str:
    return hmac.new(SECRET, f"HS256.{payload}".encode(), hashlib.sha256).hexdigest()[:32]


def issue(sub: str, jti: str, exp: int) -> str:
    payload = f"sub={sub};jti={jti};exp={exp}"
    return f"HS256.{payload}.{sign(payload)}"


def field(payload: str, key: str) -> str:
    for kv in payload.split(";"):
        if kv.startswith(key + "="):
            return kv[len(key) + 1:]
    return ""


def verify(token: str, now: int, revoked: set[str]) -> tuple[str, str]:
    parts = token.split(".")
    if len(parts) != 3 or parts[0] != "HS256":
        return "", "bad format or alg"
    if not hmac.compare_digest(parts[2], sign(parts[1])):
        return "", "signature mismatch"
    if now >= int(field(parts[1], "exp")):
        return "", "token expired"
    if field(parts[1], "jti") in revoked:
        return "", "token revoked"
    return field(parts[1], "sub"), ""


def main() -> None:
    now = 1000
    revoked: set[str] = set()

    print("== 1. 签发 ==")
    access = issue("u1", "jti-a1", now + 60)
    refresh = issue("u1", "jti-r1", now + 86400)
    print("access.jti = jti-a1，exp = +60s")
    print("refresh.jti = jti-r1，exp = +86400s")

    print("== 2. 校验正常 token ==")
    sub, err = verify(access, now, revoked)
    print(f"ok = {err == ''}，sub = {sub}")

    print("== 3. 篡改 payload ==")
    sub, err = verify(access.replace("sub=u1", "sub=u2"), now, revoked)
    print(f"ok = {err == ''}，reason = {err}")

    print("== 4. 刷新：轮换 refresh token ==")
    revoked.add("jti-r1")
    issue("u1", "jti-r2", now + 86400)
    print("新 refresh.jti = jti-r2，旧 jti-r1 进黑名单")
    sub, err = verify(refresh, now, revoked)
    print(f"用旧 refresh 再刷：ok = {err == ''}，reason = {err}")

    print("== 5. 吊销 access token ==")
    revoked.add("jti-a1")
    sub, err = verify(access, now, revoked)
    print(f"ok = {err == ''}，reason = {err}")

    print("== 6. 结论 ==")
    print("签名只保证没被篡改，吊销必须靠 jti 黑名单（或短 TTL + 版本号）")
    print("refresh 必须轮换：旧的用过一次就作废，否则泄露后可长期续命")


if __name__ == "__main__":
    main()
