package rpc

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/encrypt"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/admin"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

type RPCClient struct {
	apiConfig    *configs.APIConfig
	adminClients []admin.ServiceClient
}

func NewRPCClient(apiConfig *configs.APIConfig) (*RPCClient, error) {
	if apiConfig == nil {
		return nil, errors.New("api config should not be nil")
	}

	adminClients := []admin.ServiceClient{}

	conns := []*grpc.ClientConn{}
	for _, endpoint := range apiConfig.RPC.Endpoints {
		conn, err := grpc.Dial(endpoint, grpc.WithInsecure())
		if err != nil {
			return nil, err
		}
		conns = append(conns, conn)
	}
	if len(conns) == 0 {
		return nil, errors.New("[RPC]no available endpoints")
	}

	// node clients
	for _, conn := range conns {
		adminClients = append(adminClients, admin.NewServiceClient(conn))
	}

	return &RPCClient{
		apiConfig:    apiConfig,
		adminClients: adminClients,
	}, nil
}

func (this *RPCClient) AdminRPC() admin.ServiceClient {
	if len(this.adminClients) > 0 {
		return this.adminClients[rands.Int(0, len(this.adminClients)-1)]
	}
	return nil
}

func (this *RPCClient) Context(adminId int) context.Context {
	ctx := context.Background()
	m := maps.Map{
		"timestamp": time.Now().Unix(),
		"adminId":   adminId,
	}
	method, err := encrypt.NewMethodInstance(teaconst.EncryptMethod, this.apiConfig.Secret, this.apiConfig.NodeId)
	if err != nil {
		utils.PrintError(err)
		return context.Background()
	}
	data, err := method.Encrypt(m.AsJSON())
	if err != nil {
		utils.PrintError(err)
		return context.Background()
	}
	token := base64.StdEncoding.EncodeToString(data)
	ctx = metadata.AppendToOutgoingContext(ctx, "nodeId", this.apiConfig.NodeId, "token", token)
	return ctx
}
