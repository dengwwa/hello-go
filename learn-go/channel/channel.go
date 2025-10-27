package main

func main() {
	ch := make(chan int)

	go func() {
		// 发送操作：将数据 42 发送到 channel ch 中
		ch <- 42
	}()

	// 接收操作：从 channel ch 中接收一个值
	value := <-ch

	println(value)

	// 关闭操作：关闭 channel
	close(ch)
}
