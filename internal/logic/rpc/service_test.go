package rpc

import (
	"context"
	etcdclient "github.com/rpcxio/rpcx-etcd/client"
	"github.com/smallnest/rpcx/client"
	"ppim/internal/chatlib"
	"testing"
)

func Test_Auth(t *testing.T) {
	req := &chatlib.AuthReq{
		Uid: "123",
		Did: "bbc",
	}
	resp := &chatlib.AuthResp{}

	//discovery, _ := client.NewPeer2PeerDiscovery("tcp@localhost:9090", "")
	discovery, _ := etcdclient.NewEtcdV3Discovery("rpcx", "logic", []string{"localhost:2379"}, false, nil)

	xclient := client.NewXClient("logic", client.Failtry, client.RandomSelect, discovery, client.DefaultOption)
	defer xclient.Close()

	err := xclient.Call(context.Background(), "Auth", req, resp)
	if err != nil {
		t.Fatalf("failed to call: %v", err)
	}
	t.Logf("%v", resp.Code)
}
