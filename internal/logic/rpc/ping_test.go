package rpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"ppim/api/logic"
	"testing"
)

func Test_pingServer_Ping(t *testing.T) {
	// 连接到 gRPC 服务
	conn, err := grpc.NewClient(":9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := logic.NewPingClient(conn)

	resp, err := client.Ping(context.Background(), &logic.PingReq{})
	if err != nil {
		log.Fatalf("调用失败: %v", err)
		return
	}
	fmt.Println(resp.Msg)

}
