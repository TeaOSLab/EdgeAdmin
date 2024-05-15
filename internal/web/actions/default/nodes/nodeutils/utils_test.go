package nodeutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	_ "github.com/iwind/TeaGo/bootstrap"
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestSendMessageToCluster(t *testing.T) {
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		t.Fatal(err)
	}
	ctx := rpcClient.Context(1)

	results, err := SendMessageToCluster(ctx, 1, "test", nil, 30, false)
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(results, t)
}
