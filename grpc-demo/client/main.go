package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "grpc-demo/proto" // 引入自定义的proto包
)

// main 是程序的入口点
func main() {
	// 使用insecure连接到gRPC服务器，对于生产环境，应该使用安全的连接
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	// 创建一个Greeter服务的客户端
	client := pb.NewGreeterClient(conn)

	// 默认问候对象为"世界"，可以通过命令行参数覆盖
	name := "世界"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	// 创建一个带有超时的context，防止请求无限期等待
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 调用SayHello方法并处理响应
	res, err := client.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("请求失败: %v", err)
	}

	// 输出接收到的响应消息
	log.Printf("收到响应: %s", res.GetMessage())
}
