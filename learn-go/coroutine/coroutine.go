package main

import (
	"fmt"
	"time"
)

func main() {
	// 协程

	// 大部分情况下什么都不会输出，协程是并发执行的，系统创建协程需要时间，而在此之前，主协程早已运行结束，一旦主线程退出，其他子协程也就自然退出了
	/*	go fmt.Println("hello world, g1!")
		go hellp()
		go func() {
			fmt.Println("hello world，g3!")
		}()*/

	// 这是一个在循环体中开启协程的例子，永远也无法精准的预判到它到底会输出什么。可能子协程还没开始运行，主协程就已经结束了
	/*	fmt.Println("start")
		for i := 0; i < 10; i++ {
			go fmt.Println(i)
			//time.Sleep(time.Millisecond)
		}
		//time.Sleep(time.Millisecond)
		fmt.Println("end")*/

	// go 中常用的控制并发的手段：channel、WaitGroup、Context

	/*
		channel
	*/
	// 管道声明语句，此时管道未初始化，其值为 nil，不可以直接使用
	var ch0 chan int
	fmt.Println("管道声明但未初始化：", ch0)

	//创建缓冲区的管道
	// 没有设置缓冲区的管道，默认是阻塞的
	intCh0 := make(chan int)
	// 设置缓冲区为1的管道，
	intCh := make(chan int, 1)

	// 关闭管道
	defer close(intCh0)
	defer close(intCh)

	// 读写管道：管道中的数据流动方式与队列一样，即先进先出（FIFO），协程对于管道的操作是同步的，在某一个时刻，只有一个协程能够对其写入数据，同时也只有一个协程能够读取管道中的数据。
	// ch <- x :表示写入管道，如果管道已满，则阻塞，直到有数据从管道中取出
	fmt.Println("写入管道：")
	intCh <- 1
	// <- ch :表示读取管道，如果管道为空，则阻塞，直到有数据写入管道
	fmt.Println("读取管道：")
	value, ok := <-intCh
	if ok {
		fmt.Println("管道值：", value)
	}

	fmt.Println()
	// 无缓冲区的管道:
	//对于无缓冲管道而言，因为缓冲区容量为 0，所以不会临时存放任何数据。正因为无缓冲管道无法存放数据，
	//在向管道写入数据时必须立刻有其他协程来读取数据，否则就会阻塞等待，读取数据时也是同理，这也解释了为什么下面看起来很正常的代码会发生死锁
	fmt.Println("无缓冲区管道")
	intCh1 := make(chan int)
	defer close(intCh1)
	fmt.Println("写取无缓冲区管道：")
	// 同步执行无缓冲区管道的读写会造成 deadlock
	//intCh1 <- 123

	//无缓冲管道不应该同步的使用，正确来说应该开启一个新的协程来发送数据
	go func() {
		intCh1 <- 123
	}()
	n := <-intCh1
	fmt.Println("读无缓冲区管道值：", n)

	fmt.Println()

	// 带缓冲区的管道:
	// 带缓冲区的管道，其容量为 1，所以可以临时存放一个数据。
	fmt.Println("带缓冲区管道：")
	intCh2 := make(chan int, 1)
	defer close(intCh2)
	fmt.Println("写入带缓冲区管道：")
	// 带缓冲区管道可以同步使用，但这种同步读写的方式是非常危险的，一旦管道缓冲区空了或者满了，将会永远阻塞下去，因为没有其他协程来向管道中写入或读取数据
	intCh2 <- 123
	fmt.Println("读带缓冲区管道值：", <-intCh2)

	fmt.Println()
	ch := make(chan int, 5)
	chW := make(chan struct{})
	chR := make(chan struct{})
	defer func() {
		close(ch)
		close(chW)
		close(chR)
	}()

	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
			fmt.Println("写入管道：", i)

		}
		chW <- struct{}{}
	}()

	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond)
			fmt.Println("读取管道：", <-ch)
		}
		chR <- struct{}{}
	}()
	fmt.Println("length", len(ch), cap(ch))
	fmt.Println("写入完成...", <-chW)
	fmt.Println("读取完成...", <-chR)
}

func hellp() {
	fmt.Println("hello world，g2!")
}
