use std::future::Future;
use std::pin::Pin;
use std::sync::atomic::{AtomicBool, Ordering};
use std::sync::Arc;
use std::task::{Context, Poll, Wake, Waker};
use std::thread;
use std::time::{Duration, Instant};

// 把当前线程的 Waker 桥接成 park/unpark
struct ThreadWaker(thread::Thread);

impl Wake for ThreadWaker {
    fn wake(self: Arc<Self>) {
        self.0.unpark();
    }
}

// 极简 executor：poll -> Pending 就 park，被 wake 后继续 poll
fn block_on<F: Future>(fut: F) -> F::Output {
    let mut fut = Box::pin(fut);
    let waker = Waker::from(Arc::new(ThreadWaker(thread::current())));
    let mut cx = Context::from_waker(&waker);
    loop {
        match fut.as_mut().poll(&mut cx) {
            Poll::Ready(v) => return v,
            Poll::Pending => thread::park(),
        }
    }
}

// 第一次 poll 返回 Pending 并自我唤醒，第二次 Ready
struct YieldOnce(bool);

impl Future for YieldOnce {
    type Output = &'static str;

    fn poll(mut self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Self::Output> {
        if self.0 {
            Poll::Ready("让出后恢复")
        } else {
            self.0 = true;
            cx.waker().wake_by_ref();
            Poll::Pending
        }
    }
}

// 延时 Future：第一次 poll 起一个线程，睡完置位并唤醒
struct Delay {
    done: Arc<AtomicBool>,
}

impl Future for Delay {
    type Output = ();

    fn poll(self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<()> {
        if self.done.load(Ordering::SeqCst) {
            return Poll::Ready(());
        }
        let done = Arc::clone(&self.done);
        let waker = cx.waker().clone();
        thread::spawn(move || {
            thread::sleep(Duration::from_millis(30));
            done.store(true, Ordering::SeqCst);
            waker.wake();
        });
        Poll::Pending
    }
}

// 超时包装：每次 poll 先看 deadline，过了就直接 Err，内部 Future 不再被 poll
struct Timeout<F> {
    fut: Pin<Box<F>>,
    deadline: Instant,
}

impl<F: Future> Future for Timeout<F> {
    type Output = Result<F::Output, &'static str>;

    fn poll(mut self: Pin<&mut Self>, cx: &mut Context<'_>) -> Poll<Self::Output> {
        if Instant::now() >= self.deadline {
            return Poll::Ready(Err("timeout"));
        }
        match self.fut.as_mut().poll(cx) {
            Poll::Ready(v) => Poll::Ready(Ok(v)),
            Poll::Pending => Poll::Pending,
        }
    }
}

fn main() {
    let v = block_on(async {
        let a = YieldOnce(false).await;
        let b = YieldOnce(false).await;
        format!("{} + {}", a, b)
    });
    println!("block_on 结果 = {}", v);

    // 超时：Delay 睡 30ms，deadline 只有 1ms，第二次 poll 直接 Err
    let slow = Timeout {
        fut: Box::pin(Delay {
            done: Arc::new(AtomicBool::new(false)),
        }),
        deadline: Instant::now() + Duration::from_millis(1),
    };
    println!("超时结果 = {:?}", block_on(slow));

    // 未超时：立即 Ready 的 async 块
    let fast = Timeout {
        fut: Box::pin(async { 42 }),
        deadline: Instant::now() + Duration::from_millis(100),
    };
    println!("未超时结果 = {:?}", block_on(fast));

    println!("生产环境用 tokio + tokio::time::timeout / CancellationToken 做同样的事");
}
