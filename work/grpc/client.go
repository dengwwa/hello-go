package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	status2 "google.golang.org/grpc/status"
	resource "learn-go/work/grpc/api"
	"log"
)

func main() {
	opts := []grpc.DialOption{
		// 使用不安全连接，生产环境应使用TLS
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// 建立gRPC连接
	conn, err := grpc.Dial("10.30.60.46:8848", opts...)
	if err != nil {
		fmt.Println("did not connect:", err)
		return
	}

	// 创建gRPC客户端
	client := resource.NewResourceClient(conn)

	appIDs := []uint64{
		//12,
		//13,
		//14,
		//15,
		//16,
		//17,
		//30,
		//39,
		//42,
		//44,
		//52,
		//53,
		//54,
		//60,
		//64,
		//66,
		//70,
		//72,
		//73,
		//74,
		//76,
		//77,
		//80,
		//81,
		//82,
		//83,
		//84,
		//85,
		//86,
		//87,
		//88,
		//89,
		//90,
		//93,
		//94,
		//95,
		//97,
		//99,
		//100,
		//106,
		//110,
		//114,
		//115,
		//116,
		//117,
		//118,
		//121,
		//122,
		//123,
		//124,
	}

	for _, id := range appIDs {
		_, err = client.Apply(context.Background(), &resource.ApplyRequest{AppId: id})
		if status, _ := status2.FromError(err); status != nil {
			log.Fatalf("error: %v", status)
		} else {
			fmt.Printf("success: %v\n", id)
		}
	}

}
