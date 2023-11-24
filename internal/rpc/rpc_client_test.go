package rpc

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	_ "github.com/iwind/TeaGo/bootstrap"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"testing"
	"time"
)

func TestRPCClient_NodeRPC(t *testing.T) {
	before := time.Now()
	defer func() {
		t.Log(time.Since(before).Seconds()*1000, "ms")
	}()
	config, err := configs.LoadAPIConfig()
	if err != nil {
		t.Fatal(err)
	}
	rpc, err := NewRPCClient(config, true)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := rpc.AdminRPC().LoginAdmin(rpc.Context(0), &pb.LoginAdminRequest{
		Username: "admin",
		Password: stringutil.Md5("123456"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestRPC_Dial_HTTP(t *testing.T) {
	client, err := NewRPCClient(&configs.APIConfig{
		RPCEndpoints: []string{"https://127.0.0.1:8003"},
		NodeId:       "a7e55782dab39bce0901058a1e14a0e6",
		Secret:       "lvyPobI3BszkJopz5nPTocOs0OLkEJ7y",
	}, true)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.NodeRPC().FindEnabledNode(client.Context(1), &pb.FindEnabledNodeRequest{NodeId: 4})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Node)
}

func TestRPC_Dial_HTTP_2(t *testing.T) {
	client, err := NewRPCClient(&configs.APIConfig{
		RPCEndpoints: []string{"https://127.0.0.1:8003"},
		NodeId:       "a7e55782dab39bce0901058a1e14a0e6",
		Secret:       "lvyPobI3BszkJopz5nPTocOs0OLkEJ7y",
	}, true)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.NodeRPC().FindEnabledNode(client.Context(1), &pb.FindEnabledNodeRequest{NodeId: 4})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Node)
}

func TestRPC_Dial_HTTPS(t *testing.T) {
	client, err := NewRPCClient(&configs.APIConfig{
		RPCEndpoints: []string{"https://127.0.0.1:8004"},
		NodeId:       "a7e55782dab39bce0901058a1e14a0e6",
		Secret:       "lvyPobI3BszkJopz5nPTocOs0OLkEJ7y",
	}, true)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.NodeRPC().FindEnabledNode(client.Context(1), &pb.FindEnabledNodeRequest{NodeId: 4})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Node)
}

func BenchmarkNewRPCClient(b *testing.B) {
	config, err := configs.LoadAPIConfig()
	if err != nil {
		b.Fatal(err)
	}
	rpc, err := NewRPCClient(config, true)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := rpc.AdminRPC().LoginAdmin(rpc.Context(0), &pb.LoginAdminRequest{
			Username: "admin",
			Password: stringutil.Md5("123456"),
		})
		if err != nil {
			b.Fatal(err)
		}
		_ = resp
	}
}

func BenchmarkNewRPCClient_2(b *testing.B) {
	config, err := configs.LoadAPIConfig()
	if err != nil {
		b.Fatal(err)
	}
	rpc, err := NewRPCClient(config, true)
	if err != nil {
		b.Fatal(err)
	}

	var conn = rpc.AdminRPC()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp, err := conn.LoginAdmin(rpc.Context(0), &pb.LoginAdminRequest{
			Username: "admin",
			Password: stringutil.Md5("123456"),
		})
		if err != nil {
			b.Fatal(err)
		}
		_ = resp
	}
}
