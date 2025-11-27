package main

import (
	"encoding/json"
	"fmt"
)

type User struct {
	Name  string // 会被序列化
	Age   int    // 会被序列化
	Email string `json:"-"`               // 导出字段，但会被 JSON 忽略
	Phone string `json:"phone,omitempty"` // 空值时忽略（额外知识点）
}

func main() {
	user := User{
		Name:  "Bob",
		Age:   30,
		Email: "bob@example.com", // 该字段不会出现在 JSON 中
		Phone: "1",               // 空值，会被 omitempty 忽略
	}

	data, _ := json.Marshal(user)
	fmt.Println(string(data)) // 输出: {"name":"Bob","age":30}
}
