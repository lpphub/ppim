package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"ppim/api/logic"
	"testing"
)

func Test_Auth(t *testing.T) {
	// 连接到 gRPC 服务
	conn, err := grpc.NewClient(":9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := logic.NewLogicClient(conn)

	resp, err := client.Auth(context.Background(), &logic.AuthReq{
		Uid:   "456",
		Did:   "pc01",
		Token: "bbb",
	})
	if err != nil {
		log.Fatalf("调用失败: %v", err)
		return
	}
	fmt.Printf("%v\n", resp)

}
