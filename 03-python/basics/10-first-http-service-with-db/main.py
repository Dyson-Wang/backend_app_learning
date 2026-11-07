import json
import os
import signal
import threading
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
from typing import Any
from urllib.parse import urlparse

# 线程安全的数据层：ThreadingHTTPServer 每个连接一个线程，必须加锁
_lock = threading.Lock()
_users: dict[int, dict[str, Any]] = {}
_next_id = 1


class Handler(BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"

    def log_message(self, fmt: str, *args: Any) -> None:
        print(f"[access] {self.address_string()} {fmt % args}")

    def _send(self, status: int, payload: dict[str, Any]) -> None:
        body = json.dumps(payload, ensure_ascii=False).encode("utf-8")
        self.send_response(status)
        self.send_header("Content-Type", "application/json; charset=utf-8")
        self.send_header("Content-Length", str(len(body)))
        self.send_header("X-Request-ID", self.headers.get("X-Request-ID", "-"))
        self.end_headers()
        self.wfile.write(body)

    def _fail(self, status: int, code: str, message: str) -> None:
        self._send(status, {"error": {"code": code, "message": message}})

    def do_GET(self) -> None:
        path = urlparse(self.path).path
        if path == "/healthz":
            return self._send(200, {"status": "ok"})
        if path == "/readyz":
            return self._send(200, {"status": "ready"})
        if path == "/users":
            with _lock:
                return self._send(200, {"data": list(_users.values())})

        parts = path.strip("/").split("/")
        if len(parts) == 2 and parts[0] == "users" and parts[1].isdigit():
            with _lock:
                user = _users.get(int(parts[1]))
            return self._send(200, {"data": user}) if user else self._fail(404, "NOT_FOUND", "user not found")
        return self._fail(404, "NOT_FOUND", "route not found")

    def do_POST(self) -> None:
        path = urlparse(self.path).path
        if path != "/users":
            return self._fail(404, "NOT_FOUND", "route not found")
        try:
            length = int(self.headers.get("Content-Length", "0"))
            if length > 1_000_000:
                return self._fail(413, "TOO_LARGE", "body too large")
            raw = self.rfile.read(length).decode("utf-8")
            name = str(json.loads(raw or "{}").get("name", "")).strip()
        except (ValueError, json.JSONDecodeError):
            return self._fail(400, "BAD_JSON", "invalid JSON body")
        if not name:
            return self._fail(400, "VALIDATION", "name is required")

        global _next_id
        with _lock:
            user = {"id": _next_id, "name": name}
            _users[user["id"]] = user
            _next_id += 1
        return self._send(201, {"data": user})


def main() -> None:
    port = int(os.environ.get("PORT", "8082"))
    server = ThreadingHTTPServer(("0.0.0.0", port), Handler)
    server.daemon_threads = True

    # server.shutdown() 必须在另一个线程调用，否则会死锁
    def stop(*_: Any) -> None:
        print("signal received, draining")
        threading.Thread(target=server.shutdown, daemon=True).start()

    signal.signal(signal.SIGTERM, stop)
    signal.signal(signal.SIGINT, stop)

    print(f"listening on http://localhost:{port}", flush=True)
    try:
        server.serve_forever()
    finally:
        server.server_close()
        print("stopped")


if __name__ == "__main__":
    main()
