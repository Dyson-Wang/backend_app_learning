use std::thread;
use std::time::Duration;

#[derive(Debug)]
struct DownstreamError {
    message: String,
    retryable: bool,
}

#[derive(Debug)]
enum Outcome {
    Success { attempt: u32, value: String },
    Permanent { attempt: u32, error: String },
    Exhausted { attempt: u32, error: String },
}

fn do_with_retry<F>(max_attempts: u32, base_ms: u64, mut f: F) -> Outcome
where
    F: FnMut(u32) -> Result<String, DownstreamError>,
{
    let mut last_error = String::new();

    for attempt in 1..=max_attempts {
        match f(attempt) {
            Ok(value) => return Outcome::Success { attempt, value },
            Err(e) => {
                if !e.retryable {
                    return Outcome::Permanent {
                        attempt,
                        error: e.message,
                    };
                }
                last_error = e.message;
                if attempt == max_attempts {
                    break;
                }
                let delay = base_ms * (1u64 << (attempt - 1));
                println!(
                    "  attempt {} 失败（{}），{}ms 后重试",
                    attempt, last_error, delay
                );
                thread::sleep(Duration::from_millis(delay));
            }
        }
    }

    Outcome::Exhausted {
        attempt: max_attempts,
        error: last_error,
    }
}

fn main() {
    println!("== 1. 可重试错误：指数退避后成功 ==");
    let mut calls = 0;
    let r1 = do_with_retry(4, 10, |_| {
        calls += 1;
        if calls < 4 {
            Err(DownstreamError {
                message: "503 upstream".to_string(),
                retryable: true,
            })
        } else {
            Ok(format!("第 {} 次成功", calls))
        }
    });
    println!("结果: {:?}", r1);

    println!("== 2. 永久错误：立即停止 ==");
    let r2 = do_with_retry(4, 10, |_| {
        Err(DownstreamError {
            message: "400 bad request".to_string(),
            retryable: false,
        })
    });
    println!("结果: {:?}", r2);

    println!("== 3. 一直失败：耗尽 3 次 ==");
    let r3 = do_with_retry(3, 5, |_| {
        Err(DownstreamError {
            message: "timeout".to_string(),
            retryable: true,
        })
    });
    println!("结果: {:?}", r3);

    println!("结论：只重试可重试错误，退避必须指数增长；生产还要加抖动，避免重试风暴");
}
