package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func Test1(t *testing.T) {
	// 大部分情况下什么都不会输出，协程是并发执行的，系统创建协程需要时间，而在此之前，主协程早已运行结束，一旦主线程退出，其他子协程也就自然退出了
	go fmt.Println("hello world, g1!")
	go func() {
		fmt.Println("hello world，g2!")
	}()

	go func() {
		fmt.Println("hello world，g3!")
	}()
}

func Test2(t *testing.T) {
	// 这是一个在循环体中开启协程的例子，永远也无法精准的预判到它到底会输出什么。可能子协程还没开始运行，主协程就已经结束了
	fmt.Println("start")
	for i := 0; i < 10; i++ {
		go fmt.Println(i)
		time.Sleep(time.Millisecond) // 完整且有序输出
	}
	time.Sleep(time.Millisecond) // 完整输出但是无须
	fmt.Println("end")
}

// go 中常用的控制并发的手段：channel、WaitGroup、Context

func TestChannel1(t *testing.T) {
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
}

func TestChannel2(t *testing.T) {
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
	// 带缓冲区管道可以同步使用，但这种同步读写的方式是非常危险的，一旦管道缓冲区空了或者满了，
	//将会永远阻塞下去，因为没有其他协程来向管道中写入或读取数据
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
	fmt.Println("写入完成...", <-chW)
	fmt.Println("读取完成...", <-chR)
}

/*
  - 读写无缓冲管道：当对一个无缓冲管道直接进行同步读写操作都会导致该协程阻塞
  - 读取空缓冲区的管道：当读取一个缓冲区为空的管道时，会导致该协程阻塞
  - 写入满缓冲区的管道：当管道的缓冲区已满，对其写入数据会导致该协程阻塞
  - 管道为 nil：当管道为 nil 时，无论怎样读写都会导致当前协程阻塞
  - panic情况：
    1.关闭一个 nil 管道：当管道为 nil 时，使用 close 函数对其进行关闭操作会导致 panic`
    2.写入已关闭的管道：对一个已关闭的管道写入数据会导致 panic
    3.关闭已关闭的管道：在一些情况中，管道可能经过层层传递，调用者或许也不知道到底该由谁来关闭管道，如此一来，可能会发生关闭一个已经关闭了的管道，就会发生 panic。
  - 单向管道：
*/
func TestChannel3(t *testing.T) {
	//	利用管道的阻塞条件，可以很轻易的写出一个主协程等待子协程执行完毕的例子
	ch3 := make(chan struct{})
	defer close(ch3)
	go func() {
		fmt.Println(2)
		ch3 <- struct{}{}
	}()
	<-ch3
	fmt.Println(1)

	fmt.Println()
	// 基于缓冲区管道实现一个简单的互斥锁
	var count = 0
	var lock = make(chan struct{}, 1)
	go Add(lock, &count)
	go Sub(lock, &count)
	time.Sleep(time.Millisecond)
}

func Add(lock chan struct{}, count *int) {
	// lock
	lock <- struct{}{}
	fmt.Println("当前计数为", count, "执行加法")
	*count++
	// unlock
	<-lock
}

func Sub(lock chan struct{}, count *int) {
	lock <- struct{}{}
	fmt.Println("当前计数为", count, "执行减法")
	*count--
	<-lock
}

func TestChannel4(t *testing.T) {
	ch := make(chan int, 10)
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()
	for i := range ch {
		fmt.Println(i)
	}
}

/**
 * WaitGroup：sync.WaitGroup 是 sync 包下提供的一个结构体，WaitGroup 即等待执行，
 * 使用它可以很轻易的实现等待一组协程的效果。该结构体只对外暴露三个方法。
 *1. func (wg *WaitGroup) Add(delta int)：该方法用于增加等待的协程数，delta 表示要增加的协程数。
 *2. func (wg *WaitGroup) Done()：该方法用于减少等待的协程数，当调用该方法时，Wait 方法会返回。
 *3. func (wg *WaitGroup) Wait()：该方法用于等待一组协程执行完毕。
 */
func TestWaitGroup1(t *testing.T) {
	/**
	 * WaitGroup 使用起来十分简单，属于开箱即用。
	 * 其内部的实现是计数器+信号量，程序开始时调用 Add 初始化计数，
	 * 每当一个协程执行完毕时调用 Done，计数就-1，直到减为 0，
	 * 而在此期间，主协程调用 Wait 会一直阻塞直到全部计数减为 0，然后才会被唤醒
	 */
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		fmt.Println("协程执行完毕")
		wait.Done()
	}()
	wait.Wait()
	fmt.Println("main执行完毕")
}

func TestWaitGroup2(t *testing.T) {
	var mainWait sync.WaitGroup
	var wait sync.WaitGroup
	mainWait.Add(10)
	fmt.Println("start")
	for i := 0; i < 10; i++ {
		wait.Add(1)
		go func() {
			fmt.Println(i)
			wait.Done()
			mainWait.Done()
		}()
		wait.Wait()
	}
	mainWait.Wait()
	fmt.Println("end")
}

/*
WaitGroup 通常适用于可动态调整协程数量的时候，例如事先知晓协程的数量，又或者在运行过程中需要动态调整。
WaitGroup 的值不应该被复制，复制后的值也不应该继续使用，尤其是将其作为函数参数传递时，因该传递指针而不是值。
倘若使用复制的值，计数完全无法作用到真正的 WaitGroup 上，这可能会导致主协程一直阻塞等待，程序将无法正常运行
*/
func TestWaitGroup3(t *testing.T) {
	var mainWait sync.WaitGroup
	mainWait.Add(1)
	//错误提示所有的协程都已经退出，但主协程依旧在等待，这就形成了死锁，
	//因为 hello 函数内部对一个形参 WaitGroup 调用 Done 并不会作用到原来的 mainWait 上，
	//所以应该使用指针来进行传递。
	hello(&mainWait)
	mainWait.Wait()
	fmt.Println("main执行完毕")
}

func hello(wait *sync.WaitGroup) {
	fmt.Println("hello")
	wait.Done()
}

/*
Context 译为上下文，是 Go 提供的一种并发控制的解决方案，相比于管道和 WaitGroup，它可以更好的控制子孙协程以及层级更深的协程。Context 本身是一个接口，
只要实现了该接口都可以称之为上下文例如著名 Web 框架 Gin 中的 gin.Context。context 标准库也提供了几个实现，分别是：
  - emptyCtx
  - cancelCtx
  - timerCtx
  - valueCtx
*/
func TestContext(t *testing.T) {

}
