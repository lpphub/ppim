package rpc

import (
	"context"
	"github.com/smallnest/rpcx/client"
	"testing"
)

func Test_Auth(t *testing.T) {
	req := &AuthReq{
		Uid: "123",
		Did: "bbc",
	}
	resp := &AuthResp{}

	discovery, _ := client.NewPeer2PeerDiscovery("tcp@localhost:9090", "")
	xclient := client.NewXClient("logic", client.Failtry, client.RandomSelect, discovery, client.DefaultOption)
	defer xclient.Close()

	err := xclient.Call(context.Background(), "Auth", req, resp)
	if err != nil {
		t.Fatalf("failed to call: %v", err)
	}
	t.Logf("%v", resp.Code)
}
