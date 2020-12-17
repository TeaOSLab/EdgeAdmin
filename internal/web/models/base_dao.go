package models

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/iwind/TeaGo/logs"
)

type BaseDAO struct {
}

func (this *BaseDAO) RPC() *rpc.RPCClient {
	client, err := rpc.SharedRPC()
	if err != nil {
		logs.Println("[MODEL]get shared rpc client failed: " + err.Error())
		return nil
	}
	return client
}
