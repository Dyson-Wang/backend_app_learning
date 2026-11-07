package main

import (
	"fmt"
	"reflect"
)

//go:generate echo "占位：真实项目用 stringer/mockgen 等按 //go:generate 指令生成代码"

type User struct {
	ID   int    `json:"id" validate:"required"`
	Name string `json:"name" validate:"min=2"`
	age  int    // 未导出：反射可读类型，但不可 Set
}

func (u User) Greet(prefix string) string { return prefix + u.Name }

func main() {
	u := User{ID: 1, Name: "Alice", age: 18}

	// 1. 运行时类型信息
	t := reflect.TypeOf(u)
	fmt.Printf("类型信息: Name = %q, Kind = %s, 字段数 = %d\n", t.Name(), t.Kind(), t.NumField())

	// 2. 遍历字段 + 读取 struct tag（Go 的「注解」就是标签字符串）
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Printf("字段 %d: %s(%v) json=%q validate=%q 导出=%t\n",
			i, f.Name, f.Type, f.Tag.Get("json"), f.Tag.Get("validate"), f.PkgPath == "")
	}

	// 3. 按名字访问 / 修改字段：必须传指针，Value 才 CanSet
	pv := reflect.ValueOf(&u).Elem()
	nameField := pv.FieldByName("Name")
	fmt.Printf("按名访问: Name = %q, CanSet = %t\n", nameField.String(), nameField.CanSet())
	if nameField.CanSet() {
		nameField.SetString("Bob")
	}
	fmt.Printf("反射设值后: u.Name = %q\n", u.Name)

	// 4. 边界：未导出字段不可设置（强设会 panic）
	ageField := pv.FieldByName("age")
	fmt.Printf("反射边界: age.CanSet = %t（未导出字段不可写）\n", ageField.CanSet())

	// 5. 动态调用方法：MethodByName + Call，参数/返回都是 []reflect.Value
	method := reflect.ValueOf(u).MethodByName("Greet")
	ret := method.Call([]reflect.Value{reflect.ValueOf("你好, ")})
	fmt.Printf("动态调用: Greet 返回 = %q\n", ret[0].String())

	// 6. 反射创建值
	iv := reflect.New(reflect.TypeOf(0)).Elem()
	iv.SetInt(7)
	fmt.Printf("反射创建: new(int) 设为 %d, Interface() = %v\n", iv.Int(), iv.Interface())

	// 7. interface 的动态类型 vs 静态类型
	var boxed any = User{ID: 9}
	fmt.Printf("动态类型: 静态 any, 实际 = %v\n", reflect.TypeOf(boxed))
}
