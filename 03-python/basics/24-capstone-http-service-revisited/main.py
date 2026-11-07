import json
import threading
import time
import urllib.error
import urllib.request
from dataclasses import dataclass, field
from datetime import datetime, timedelta, timezone
from enum import Enum
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from typing import Any, Callable
from urllib.parse import urlparse

# 枚举当错误码（13 课）
class ErrCode(Enum):
    NOT_FOUND = "NOT_FOUND"
    VALIDATION = "VALIDATION"
    BAD_JSON = "BAD_JSON"


# 泛型结果容器（14 课 PEP 695）
class Result[T]:
    def __init__(self, ok: bool, data: T, status: int, code: str = "") -> None:
        self.ok = ok
        self.data = data
        self.status = status
        self.code = code

    def as_dict(self) -> dict[str, Any]:
        return {"ok": self.ok, "data": self.data, "code": self.code}


@dataclass
class User:
    uid: int
    name: str
    tags: list[str] = field(default_factory=list)


# 内存数据层 + 并发保护（09/19 课）
_lock = threading.Lock()
_users: dict[int, User] = {}
_next_id = 1
_req_seq = 0
_logs: list[str] = []


def make_logger(prefix: str) -> Callable[[str], str]:
    # 闭包（16 课）：捕获 prefix，返回专用日志函数
    def write(msg: str) -> str:
        global _req_seq
        with _lock:
            _req_seq += 1
            seq = _req_seq
        # 时间特性（18 课）：用固定基准时间演示格式化，不打印当前绝对时间
        stamp = datetime(2024, 1, 1, tzinfo=timezone.utc) + timedelta(seconds=seq)
        line = f"{prefix}[req#{seq} {stamp:%H:%M:%S}] {msg}"
        _logs.append(line)
        return line
    return write


log = make_logger("[access] ")


def validate_name(body: Any) -> str:
    # 校验 + 自定义错误处理（04 课）
    if not isinstance(body, dict):
        raise ValueError("请求体必须是对象")
    name = str(body.get("name", "")).strip()
    if not name:
        raise ValueError("name 不能为空")
    return name


def route(method: str, path: str, body: Any) -> Result[Any]:
    # match case 路由（13 课）
    match (method, path):
        case ("GET", "/healthz"):
            return Result(True, {"status": "ok"}, 200)
        case ("GET", "/users"):
            with _lock:
                data = [vars(u) for u in _users.values()]
            return Result(True, {"items": data}, 200)
        case ("POST", "/users"):
            try:
                name = validate_name(body)
            except ValueError as exc:
                return Result(False, {"message": str(exc)}, 400, ErrCode.VALIDATION.value)
            global _next_id
            with _lock:
                user = User(_next_id, name, ["new"])
                _users[user.uid] = user
                _next_id += 1
            return Result(True, vars(user), 201)
        case ("GET", p) if p.startswith("/users/"):
            raw = p.rsplit("/", 1)[-1]
            if raw.isdigit():
                with _lock:
                    u = _users.get(int(raw))
                if u:
                    return Result(True, vars(u), 200)
            return Result(False, {"message": "user not found"}, 404, ErrCode.NOT_FOUND.value)
        case _:
            return Result(False, {"message": "route not found"}, 404, ErrCode.NOT_FOUND.value)


class Handler(BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"

    def log_message(self, fmt: str, *args: Any) -> None:
        pass  # 屏蔽默认 stderr 日志，改用自定义日志

    def do_GET(self) -> None:
        self._dispatch("GET")

    def do_POST(self) -> None:
        self._dispatch("POST")

    def _send(self, result: Result[Any]) -> None:
        payload = json.dumps(result.as_dict(), ensure_ascii=False).encode("utf-8")
        self.send_response(result.status)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(payload)))
        self.send_header("X-Request-ID", str(self.headers.get("X-Request-ID", "-")))
        self.end_headers()
        self.wfile.write(payload)

    def _dispatch(self, method: str) -> None:
        path = urlparse(self.path).path
        body: Any = None
        if method == "POST":
            try:
                length = int(self.headers.get("Content-Length", "0"))
                raw = self.rfile.read(length).decode("utf-8")
                body = json.loads(raw or "{}")
            except (ValueError, json.JSONDecodeError):
                self._send(Result(False, {"message": "invalid JSON"}, 400, ErrCode.BAD_JSON.value))
                return
        log(f"{method} {path}")
        result = route(method, path, body)
        log(f"-> status={result.status} code={result.code}")
        self._send(result)


def request(port: int, method: str, path: str, payload: dict | None = None) -> tuple[int, dict]:
    # 主线程自请求（urllib 标准库）
    url = f"http://127.0.0.1:{port}{path}"
    data = None if payload is None else json.dumps(payload).encode("utf-8")
    req = urllib.request.Request(url, data=data, method=method,
                                 headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=2.0) as resp:
            return resp.status, json.loads(resp.read().decode("utf-8"))
    except urllib.error.HTTPError as err:
        return err.code, json.loads(err.read().decode("utf-8"))


def main() -> None:
    # 随机端口：绑定 127.0.0.1:0 由内核分配，避免占用固定端口
    server = ThreadingHTTPServer(("127.0.0.1", 0), Handler)
    server.daemon_threads = True
    port = server.server_address[1]
    t = threading.Thread(target=server.serve_forever, daemon=True)
    t.start()
    print(f"1) 内存版 HTTP 服务已启动，随机端口={port}（自请求后 3 秒内关闭）")

    started = time.perf_counter()
    print("2) GET /healthz")
    print(f"   {request(port, 'GET', '/healthz')}")
    print("3) 创建 / 校验 / 查询 / 错误码")
    print(f"   创建 alice: {request(port, 'POST', '/users', {'name': 'alice'})}")
    print(f"   空 name 校验: {request(port, 'POST', '/users', {'name': ''})}")
    print(f"   用户列表: {request(port, 'GET', '/users')}")
    print(f"   查 id=1: {request(port, 'GET', '/users/1')}")
    print(f"   查不存在: {request(port, 'GET', '/users/999')}")
    print(f"   未知路由: {request(port, 'GET', '/nope')}")

    elapsed = time.perf_counter() - started
    print(f"4) 自请求耗时>0={elapsed > 0}；耗时<3秒={elapsed < 3.0}")
    print("5) 服务访问日志（按请求序号，固定基准时间格式化）")
    for line in _logs:
        print("   " + line)

    # shutdown() 必须从另一个线程调用（此处主线程调，serve_forever 在子线程）
    server.shutdown()
    server.server_close()
    t.join(timeout=1.0)
    print(f"6) 已 shutdown，子线程结束={not t.is_alive()}")


if __name__ == "__main__":
    main()
