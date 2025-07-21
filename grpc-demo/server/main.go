package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "grpc-demo/proto" // 模块路径
)

// server 结构体实现了 pb.UnimplementedGreeterServer 接口，用于处理 gRPC 请求。
type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello 方法处理来自客户端的 HelloRequest 请求。
// 该方法接收一个上下文和一个 HelloRequest 对象，返回一个 HelloReply 对象和一个错误对象。
// 主要功能是构造回复消息，消息内容为 "你好, [请求的名字]!"。
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	// 记录收到的请求名字
	log.Printf("收到请求: %v", req.GetName())
	// 构造回复消息并返回
	return &pb.HelloReply{Message: "你好, " + req.GetName() + "!"}, nil
}

func main() {
	// 监听 TCP 端口 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		// 如果监听失败，记录错误并退出
		log.Fatalf("监听失败: %v", err)
	}

	// 创建一个新的 gRPC 服务器实例
	s := grpc.NewServer()
	// 将 server 实例注册到 gRPC 服务器中
	pb.RegisterGreeterServer(s, &server{})

	// 记录服务启动信息
	log.Println("服务端运行中: 50051 端口")
	// 使用监听器启动 gRPC 服务
	if err := s.Serve(lis); err != nil {
		// 如果服务启动失败，记录错误并退出
		log.Fatalf("服务失败: %v", err)
	}
}
