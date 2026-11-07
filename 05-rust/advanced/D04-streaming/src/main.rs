const CHUNKS: i32 = 5;
const PER_CHUNK: i32 = 20;

fn buffered() -> (i32, i32) {
    let mut t = 0;
    for _ in 0..CHUNKS {
        t += PER_CHUNK;
    }
    (t, t) // 全部生成完才开始发
}

fn streamed() -> (i32, i32) {
    let mut t = 0;
    let mut first = 0;
    for i in 0..CHUNKS {
        t += PER_CHUNK;
        if i == 0 {
            first = t; // 第一块生成完立刻发
        }
    }
    (first, t)
}

fn main() {
    println!("== 1. 缓冲模式：全部生成完再发 ==");
    let (fb, total) = buffered();
    println!("首字节 = {}ms，总耗时 = {}ms", fb, total);

    println!("== 2. 流式模式：生成一块发一块 ==");
    let (fb, total) = streamed();
    println!("首字节 = {}ms，总耗时 = {}ms", fb, total);

    println!("== 3. 对比 ==");
    println!("总耗时一样，首字节差了 4 倍：用户感知的是首字节，不是总耗时");
    println!("SSE / chunked / WebSocket 的价值就在这里；代价是不能压缩、不能算 Content-Length");

    println!("== 4. 结论 ==");
    println!("大响应先想清楚：客户端要等全量，还是能边收边渲染");
    println!("流式要处理背压：客户端读得慢，服务端不能无限缓冲");
}
