#[derive(Debug, Clone, Copy, PartialEq)]
enum Code {
    Required,
    Length,
    Format,
    Range,
}

impl Code {
    fn as_str(self) -> &'static str {
        match self {
            Code::Required => "REQUIRED",
            Code::Length => "LENGTH",
            Code::Format => "FORMAT",
            Code::Range => "RANGE",
        }
    }
}

#[derive(Debug)]
struct FieldError {
    code: Code,
    field: &'static str,
    message: String,
}

struct CreateUser {
    name: Option<String>,
    email: Option<String>,
    age: Option<i32>,
}

fn validate(u: &CreateUser) -> Vec<FieldError> {
    let mut errors = Vec::new();

    let name = u.name.as_deref().unwrap_or("").trim();
    if name.is_empty() {
        errors.push(FieldError {
            code: Code::Required,
            field: "name",
            message: "name 不能为空".to_string(),
        });
    } else {
        let n = name.chars().count();
        if n < 2 || n > 20 {
            errors.push(FieldError {
                code: Code::Length,
                field: "name",
                message: "name 长度必须在 2-20 之间".to_string(),
            });
        }
    }

    let email = u.email.as_deref().unwrap_or("");
    if email.is_empty() {
        errors.push(FieldError {
            code: Code::Required,
            field: "email",
            message: "email 不能为空".to_string(),
        });
    } else if !email.contains('@') {
        errors.push(FieldError {
            code: Code::Format,
            field: "email",
            message: "email 格式不正确".to_string(),
        });
    }

    match u.age {
        None => errors.push(FieldError {
            code: Code::Required,
            field: "age",
            message: "age 不能为空".to_string(),
        }),
        Some(a) if !(0..=150).contains(&a) => errors.push(FieldError {
            code: Code::Range,
            field: "age",
            message: "age 必须在 0-150 之间".to_string(),
        }),
        Some(_) => {}
    }

    errors
}

fn report(label: &str, u: &CreateUser) {
    let errors = validate(u);
    println!("{}：{} 个错误", label, errors.len());
    for e in &errors {
        println!("  [{}] {}: {}", e.code.as_str(), e.field, e.message);
    }
}

fn main() {
    report(
        "全空",
        &CreateUser {
            name: None,
            email: None,
            age: None,
        },
    );
    report(
        "短名 + 坏邮箱 + 超龄",
        &CreateUser {
            name: Some("a".to_string()),
            email: Some("bad".to_string()),
            age: Some(200),
        },
    );
    report(
        "合法",
        &CreateUser {
            name: Some("ada".to_string()),
            email: Some("ada@example.com".to_string()),
            age: Some(36),
        },
    );

    println!("结论：一次返回全部错误，前端可逐字段标红；code 与 message 分离便于国际化");
}
