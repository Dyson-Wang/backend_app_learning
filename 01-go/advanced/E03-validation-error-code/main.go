package main

import (
	"fmt"
	"strings"
)

type CreateUser struct {
	Name  string
	Email string
	Age   *int // nil 表示未提供，避免 0 岁被当成「合法但没填」
}

func intPtr(v int) *int { return &v }

type FieldError struct {
	Code    string `json:"code"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

// 一次收集全部错误，而不是遇到第一个就返回
func Validate(u CreateUser) []FieldError {
	errs := make([]FieldError, 0, 4)
	add := func(code, field, msg string) {
		errs = append(errs, FieldError{Code: code, Field: field, Message: msg})
	}

	if strings.TrimSpace(u.Name) == "" {
		add("REQUIRED", "name", "name 不能为空")
	} else if n := len([]rune(u.Name)); n < 2 || n > 20 {
		add("LENGTH", "name", "name 长度必须在 2-20 之间")
	}

	if u.Email == "" {
		add("REQUIRED", "email", "email 不能为空")
	} else if !strings.Contains(u.Email, "@") {
		add("FORMAT", "email", "email 格式不正确")
	}

	if u.Age == nil {
		add("REQUIRED", "age", "age 不能为空")
	} else if *u.Age < 0 || *u.Age > 150 {
		add("RANGE", "age", "age 必须在 0-150 之间")
	}
	return errs
}

func report(label string, u CreateUser) {
	errs := Validate(u)
	fmt.Printf("%s：%d 个错误\n", label, len(errs))
	for _, e := range errs {
		fmt.Printf("  [%s] %s: %s\n", e.Code, e.Field, e.Message)
	}
}

func main() {
	report("全空", CreateUser{})
	report("短名 + 坏邮箱 + 超龄", CreateUser{Name: "a", Email: "bad", Age: intPtr(200)})
	report("合法", CreateUser{Name: "ada", Email: "ada@example.com", Age: intPtr(36)})
	fmt.Println("结论：一次返回全部错误，前端可逐字段标红，而不是改一个再试一个")
}
