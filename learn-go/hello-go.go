package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 创建可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())

	// 启动 worker
	go worker(ctx)

	// 让 worker 运行一段时间
	time.Sleep(3 * time.Second)

	// 取消上下文，停止 worker
	fmt.Println("准备取消工作...")
	cancel()

	// 给 worker 时间处理取消
	time.Sleep(1 * time.Second)
	fmt.Println("主程序结束")

}

func worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			// 上下文被取消，执行清理工作
			fmt.Println("工作被取消，原因:", ctx.Err())
			return
		default:
			// 正常的工作逻辑
			fmt.Println("工作中...")
			time.Sleep(1 * time.Second)
		}
	}
}
