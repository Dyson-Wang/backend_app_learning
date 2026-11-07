"""密码哈希与 RBAC。"""
import hashlib
import hmac

ROLE_PERMS = {
    "admin": {"order:read", "order:write", "user:delete"},
    "viewer": {"order:read"},
}


def hash_password(password: str, salt: str) -> str:
    return hashlib.pbkdf2_hmac("sha256", password.encode(), salt.encode(), 1000, 16).hex()


def safe_equal(a: str, b: str) -> bool:
    # 固定时间比较，避免按字节提前返回泄露信息
    return hmac.compare_digest(a.encode(), b.encode())


def allowed(role: str, perm: str) -> bool:
    return perm in ROLE_PERMS.get(role, set())


def main() -> None:
    salt = "salt-8f2a"
    hashed = hash_password("correct-horse", salt)

    print("== 1. 加盐哈希 ==")
    print(f"salt = {salt}，hash = {hashed[:16]}...")
    print(f"同一密码 + 同一盐 结果一致 = {hashed == hash_password('correct-horse', salt)}")
    print(f"同一密码 + 不同盐 结果不同 = {hashed != hash_password('correct-horse', 'salt-other')}")

    print("== 2. 校验 ==")
    print(f"正确密码 = {safe_equal(hashed, hash_password('correct-horse', salt))}")
    print(f"错误密码 = {safe_equal(hashed, hash_password('wrong', salt))}")

    print("== 3. RBAC ==")
    print(f"admin 能 order:write = {allowed('admin', 'order:write')}")
    print(f"viewer 能 order:write = {allowed('viewer', 'order:write')}")
    print(f"viewer 能 user:delete = {allowed('viewer', 'user:delete')}")

    print("== 4. 结论 ==")
    print("密码永远不存明文，也不要可逆加密；盐必须每用户独立随机")
    print("1000 轮只是演示，生产用 bcrypt/argon2id，成本参数按机器调")


if __name__ == "__main__":
    main()
