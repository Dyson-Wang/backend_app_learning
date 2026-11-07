"""优雅关闭：停止接收新连接，等在途请求跑完。"""
import http.server
import os
import signal
import socket
import threading
import time
import urllib.error
import urllib.request


class Handler(http.server.BaseHTTPRequestHandler):
    protocol_version = "HTTP/1.1"

    def do_GET(self) -> None:
        if self.path == "/slow":
            time.sleep(0.12)          # 模拟在途慢请求
            body = b"slow done"
        else:
            body = b"ok"
        self.send_response(200)
        self.send_header("Content-Length", str(len(body)))
        self.send_header("Connection", "close")
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, *args: object) -> None:
        pass


def main() -> None:
    server = http.server.ThreadingHTTPServer(("127.0.0.1", 0), Handler)
    port = server.server_address[1]
    threading.Thread(target=server.serve_forever, daemon=True).start()

    print("== 1. 发起一个 120ms 的在途请求 ==")
    result: dict[str, str] = {}

    def call(path: str) -> None:
        try:
            with urllib.request.urlopen(f"http://127.0.0.1:{port}{path}", timeout=3) as r:
                result[path] = r.read().decode()
        except Exception as exc:            # noqa: BLE001
            result[path] = f"ERR {type(exc).__name__}"

    t = threading.Thread(target=call, args=("/slow",))
    t.start()
    time.sleep(0.03)

    print("== 2. 触发优雅关闭 ==")
    t0 = time.monotonic()
    threading.Thread(target=server.shutdown, daemon=True).start()   # 停止 accept
    t.join()
    cost = int((time.monotonic() - t0) * 1000)
    print(f"body = {result['/slow']}，关闭耗时 = {cost}ms")
    server.server_close()          # 关掉监听 socket，之后连接直接被拒

    print("== 3. 关闭后新请求 ==")
    try:
        with socket.create_connection(("127.0.0.1", port), timeout=1):
            print("新请求 -> 连接成功（不应该发生）")
    except OSError as exc:
        print(f"新请求 -> 被拒绝 {type(exc).__name__}")

    print("== 4. 结论 ==")
    print("shutdown() 停止接收新连接，serve_forever 退出前会等现有请求")
    print("生产上用 signal.signal(SIGTERM) 触发，并配超时兜底")
    print(f"（本进程 pid={os.getpid()}，真实部署时由编排系统发 SIGTERM）")


if __name__ == "__main__":
    signal.signal(signal.SIGPIPE, signal.SIG_IGN)
    main()
