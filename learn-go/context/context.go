package main

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// 什么是 Context？ Context（上下文）是 Go 语言中用于传递请求范围数据、取消信号和超时时间的标准方式。它主要在 goroutine 之间传递信息。
// Context 类型：
// 1. 空的根 Context:context.Background()
// 2. 带值的 Context:context.WithValue(parent, key, value)
// 3. 可取消的 Context:context.WithCancel(parent)
// 4. 带超时的 Context:context.WithTimeout(parent, timeout)
// 5. 带截止时间的 Context:context.WithDeadline(parent, deadline)

// 最佳实践
// 1.将 Context 作为函数的第一个参数 func myFunction(ctx context.Context, otherParam string) error
// 2.不要存储 Context 在结构体中（除非必要）
// 3.总是调用 cancel 函数以避免内存泄漏
// 4.使用 context.WithValue 时要小心，只传递请求范围的数据
// 5.检查 Context 是否被取消：if err := ctx.Err(); err != nil { return err }

func main() {
	/*	contextWithBackground()
		contextWithValue()
		contextWithCancel()
		contextWithTimeout()
		contextWithDeadline()*/

	HttpExample{}.demo()
}

// 空的根 Context
func contextWithBackground() {
	fmt.Println("=========空的根 Context=========")
	fmt.Println("contextWithBackground")
	ctx := context.Background()
	fmt.Printf("Context 类型: %T\n", ctx)

	// 或者使用 TODO（当不确定用哪种 Context 时）
	ctx2 := context.TODO()
	fmt.Printf("Context 类型: %T\n", ctx2)
}

// 带值的 Context
func contextWithValue() {
	fmt.Println("=========带值的 Context=========")

	ctx := context.Background()
	// 添加键值对到 Context 中
	ctx = context.WithValue(ctx, "userID", 12345)
	ctx = context.WithValue(ctx, "username", "Alice")

	// 从 Context 中获取值
	userID := ctx.Value("userID")
	username := ctx.Value("username")

	fmt.Printf("context: %v, UserID: %v, Username: %v\n", ctx, userID, username)
}

// 可取消的 Context
func contextWithCancel() {
	fmt.Println("=========可取消的 Context=========")

	ctx, cancelFunc := context.WithCancel(context.Background())

	go func(ctx2 context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Worker: 收到取消信号，退出")
				return
			default:
				fmt.Println("Worker: 工作中...")
				time.Sleep(1 * time.Second)
			}
		}
	}(ctx)

	time.Sleep(1 * time.Second)
	cancelFunc()
	fmt.Println("程序结束")

}

// 带超时的 Context
func contextWithTimeout() {
	fmt.Println("=========带超时的 Context=========")

	// 创建 2 秒超时的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 确保释放资源

	go func(ctx2 context.Context) {
		select {
		case <-time.After(3 * time.Second):
			fmt.Println("操作完成")
		case <-ctx2.Done():
			fmt.Println("操作被取消或超时:", ctx2.Err())
		}
	}(ctx)

	// 等待操作完成或超时
	<-ctx.Done()
	fmt.Println("主程序结束:", ctx.Err())
}

// 带截止时间的 Context
func contextWithDeadline() {
	fmt.Println("=========带截止时间的 Context=========")
	// 设置截止时间为 2 秒后
	deadline := time.Now().Add(2 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()

	go func(ctx2 context.Context) {
		select {
		case <-time.After(3 * time.Second):
			fmt.Println("任务完成")
		case <-ctx.Done():
			fmt.Println("任务未在截止时间内完成")
		}
	}(ctx)

	// 等待操作完成或超时
	<-ctx.Done()
	fmt.Println("任务状态:", ctx.Err())
}

type HttpExample struct {
}

func (h HttpExample) demo() {
	fmt.Println("=========HttpExample=========")
	http.HandleFunc("/api/data", h.apiHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return
	}
}

func (h HttpExample) apiHandler(w http.ResponseWriter, r *http.Request) {
	// 从请求中获取 Context
	ctx := r.Context()

	// 创建带超时的 Context
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 执行一些耗时操作
	result, err := h.fetchData(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		return
	}

	fmt.Println(w, "数据: %s", result)
}

func (HttpExample) fetchData(ctx context.Context) (string, error) {
	// 模拟耗时操作
	select {
	case <-time.After(3 * time.Second):
		return "获取到的数据", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
