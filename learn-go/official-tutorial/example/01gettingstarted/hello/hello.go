package main

import "fmt"
import "rsc.io/quote"

// go mod init moduleName
// go mod tidy:Go 将把 quote 模块作为依赖添加，并生成一个 go.sum 文件用于模块认证
// go run hello.go
func main() {
	fmt.Println("hello world")
	fmt.Println(quote.Go())
}
