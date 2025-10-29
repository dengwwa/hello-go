package main

import (
	"fmt"
	"time"
)

// 创建一个无缓冲的 channel
func main() {
	//unbufferedChannel()
	bufferedChannel()
}

/**
 * channel 的操作
 * 1.创建 channel
 * 2.发送操作
 * 3.接收操作
 * 4.关闭 channel
 */
func channelOperations() {
	//1.创建channel
	// 创建一个可以传递 int 类型的无缓冲 channel
	ch1 := make(chan int)
	// 创建一个可以传递 string 类型，且缓冲区大小为 10 的 channel
	ch2 := make(chan int, 10)

	// 2.发送操作：将数据 42 发送到 channel ch 中
	ch1 <- 42
	ch2 <- 43

	// 3.接收操作：从 channel ch 中接收一个值
	value := <-ch1
	fmt.Println("接收到的值:", value)

	// 关闭操作：关闭 channel
	close(ch1)
}

/**
 * 无缓冲 channel 本质上是一种【同步操作】
 * （1）在单线程中，无缓冲的 channel 在发送或接收操作时，如果 channel 缓冲区已满或已空，则操作会阻塞，导致后续的接收操作永远无法执行，最终触发死锁
 * （2）在多线程中，无缓冲的 channel 需要在不同的现场中分别实现发送和接收操作，否则可能会导致数据丢失或者数据无法被接收
 */
func unbufferedChannel() {
	fmt.Println("========无缓冲 channel========")
	ch := make(chan int)

	// 单线程中同时执行发送和接收操作会触发死锁
	//ch <- 42
	//value := <-ch
	//fmt.Println("接收到的值:", value)

	go func() {
		ch <- 42
	}()

	time.Sleep(1 * time.Second)
	// 尝试接收数据
	select {
	case value := <-ch:
		println("接收成功:", value)
	default:
		println("接收失败")
	}

	fmt.Println("========无缓冲 channel========")

}

/**
 * 有缓冲 channel 本质上是一种【异步通信】
 * 异步通信： 有缓冲 Channel 就像是“邮筒”。
 *	 a.只要缓冲区未满，发送者就可以不阻塞地继续发送数据。
 *	 b.只要缓冲区不为空，接收者就可以不阻塞地继续接收数据。
 *	 c.只有当缓冲区已满时，发送者才会阻塞；只有当缓冲区为空时，接收者才会阻塞。
 *特点： 它解耦了发送和接收的时机，提供了一定的异步能力。
 */
func bufferedChannel() {
	fmt.Println("========有缓冲 channel========")
	ch := make(chan int, 2)
	// 单线程中同事执行发送和接收操作缓冲区不满时，会正常执行
	ch <- 42
	ch <- 43
	//ch <- 43 // 如果加上这行，会阻塞，因为缓冲区已满，没有接收者,会导致死锁
	fmt.Println("接收到的值:", <-ch)
	fmt.Println("接收到的值:", <-ch)

}
