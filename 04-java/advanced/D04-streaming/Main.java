public class Main {

    static final int CHUNKS = 5;
    static final int PER_CHUNK = 20;

    public static void main(String[] args) {
        System.out.println("== 1. 缓冲模式：全部生成完再发 ==");
        int t = 0;
        for (int i = 0; i < CHUNKS; i++) {
            t += PER_CHUNK;
        }
        System.out.println("首字节 = " + t + "ms，总耗时 = " + t + "ms");

        System.out.println("== 2. 流式模式：生成一块发一块 ==");
        t = 0;
        int firstByte = 0;
        for (int i = 0; i < CHUNKS; i++) {
            t += PER_CHUNK;
            if (i == 0) {
                firstByte = t;
            }
        }
        System.out.println("首字节 = " + firstByte + "ms，总耗时 = " + t + "ms");

        System.out.println("== 3. 对比 ==");
        System.out.println("总耗时一样，首字节差了 4 倍：用户感知的是首字节，不是总耗时");
        System.out.println("SSE / chunked / WebSocket 的价值就在这里；代价是不能压缩、不能算 Content-Length");

        System.out.println("== 4. 结论 ==");
        System.out.println("大响应先想清楚：客户端要等全量，还是能边收边渲染");
        System.out.println("流式要处理背压：客户端读得慢，服务端不能无限缓冲");
    }
}
