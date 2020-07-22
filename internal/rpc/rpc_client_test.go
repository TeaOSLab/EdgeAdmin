package rpc

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/admin"
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
	rpc, err := NewRPCClient(config)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := rpc.AdminRPC().Login(rpc.Context(0), &admin.LoginRequest{
		Username: "admin",
		Password: stringutil.Md5("123456"),
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}
