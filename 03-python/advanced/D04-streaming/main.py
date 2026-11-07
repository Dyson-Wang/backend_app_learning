"""流式输出 vs 缓冲输出（逻辑时钟）。"""

CHUNKS = 5
PER_CHUNK = 20


def buffered() -> tuple[int, int]:
    t = 0
    for _ in range(CHUNKS):
        t += PER_CHUNK
    return t, t          # 全部生成完才开始发


def streamed() -> tuple[int, int]:
    t = 0
    first = 0
    for i in range(CHUNKS):
        t += PER_CHUNK
        if i == 0:
            first = t    # 第一块生成完立刻发
    return first, t


def main() -> None:
    print("== 1. 缓冲模式：全部生成完再发 ==")
    fb, total = buffered()
    print(f"首字节 = {fb}ms，总耗时 = {total}ms")

    print("== 2. 流式模式：生成一块发一块 ==")
    fb, total = streamed()
    print(f"首字节 = {fb}ms，总耗时 = {total}ms")

    print("== 3. 对比 ==")
    print("总耗时一样，首字节差了 4 倍：用户感知的是首字节，不是总耗时")
    print("SSE / chunked / WebSocket 的价值就在这里；代价是不能压缩、不能算 Content-Length")

    print("== 4. 结论 ==")
    print("大响应先想清楚：客户端要等全量，还是能边收边渲染")
    print("流式要处理背压：客户端读得慢，服务端不能无限缓冲")


if __name__ == "__main__":
    main()
