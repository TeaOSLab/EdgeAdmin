package rpc

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/encrypt"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

type RPCClient struct {
	apiConfig            *configs.APIConfig
	adminClients         []pb.AdminServiceClient
	nodeClients          []pb.NodeServiceClient
	nodeGrantClients     []pb.NodeGrantServiceClient
	nodeClusterClients   []pb.NodeClusterServiceClient
	nodeIPAddressClients []pb.NodeIPAddressServiceClient
	serverClients        []pb.ServerServiceClient
	apiNodeClients       []pb.APINodeServiceClient
}

func NewRPCClient(apiConfig *configs.APIConfig) (*RPCClient, error) {
	if apiConfig == nil {
		return nil, errors.New("api config should not be nil")
	}

	adminClients := []pb.AdminServiceClient{}
	nodeClients := []pb.NodeServiceClient{}
	nodeGrantClients := []pb.NodeGrantServiceClient{}
	nodeClusterClients := []pb.NodeClusterServiceClient{}
	nodeIPAddressClients := []pb.NodeIPAddressServiceClient{}
	serverClients := []pb.ServerServiceClient{}
	apiNodeClients := []pb.APINodeServiceClient{}

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
		adminClients = append(adminClients, pb.NewAdminServiceClient(conn))
		nodeClients = append(nodeClients, pb.NewNodeServiceClient(conn))
		nodeGrantClients = append(nodeGrantClients, pb.NewNodeGrantServiceClient(conn))
		nodeClusterClients = append(nodeClusterClients, pb.NewNodeClusterServiceClient(conn))
		nodeIPAddressClients = append(nodeIPAddressClients, pb.NewNodeIPAddressServiceClient(conn))
		serverClients = append(serverClients, pb.NewServerServiceClient(conn))
		apiNodeClients = append(apiNodeClients, pb.NewAPINodeServiceClient(conn))
	}

	return &RPCClient{
		apiConfig:            apiConfig,
		adminClients:         adminClients,
		nodeClients:          nodeClients,
		nodeGrantClients:     nodeGrantClients,
		nodeClusterClients:   nodeClusterClients,
		nodeIPAddressClients: nodeIPAddressClients,
		serverClients:        serverClients,
		apiNodeClients:       apiNodeClients,
	}, nil
}

func (this *RPCClient) AdminRPC() pb.AdminServiceClient {
	if len(this.adminClients) > 0 {
		return this.adminClients[rands.Int(0, len(this.adminClients)-1)]
	}
	return nil
}

func (this *RPCClient) NodeRPC() pb.NodeServiceClient {
	if len(this.nodeClients) > 0 {
		return this.nodeClients[rands.Int(0, len(this.nodeClients)-1)]
	}
	return nil
}

func (this *RPCClient) NodeGrantRPC() pb.NodeGrantServiceClient {
	if len(this.nodeGrantClients) > 0 {
		return this.nodeGrantClients[rands.Int(0, len(this.nodeGrantClients)-1)]
	}
	return nil
}

func (this *RPCClient) NodeClusterRPC() pb.NodeClusterServiceClient {
	if len(this.nodeClusterClients) > 0 {
		return this.nodeClusterClients[rands.Int(0, len(this.nodeClusterClients)-1)]
	}
	return nil
}

func (this *RPCClient) NodeIPAddressRPC() pb.NodeIPAddressServiceClient {
	if len(this.nodeIPAddressClients) > 0 {
		return this.nodeIPAddressClients[rands.Int(0, len(this.nodeIPAddressClients)-1)]
	}
	return nil
}

func (this *RPCClient) ServerRPC() pb.ServerServiceClient {
	if len(this.serverClients) > 0 {
		return this.serverClients[rands.Int(0, len(this.serverClients)-1)]
	}
	return nil
}

func (this *RPCClient) APINodeRPC() pb.APINodeServiceClient {
	if len(this.serverClients) > 0 {
		return this.apiNodeClients[rands.Int(0, len(this.apiNodeClients)-1)]
	}
	return nil
}

func (this *RPCClient) Context(adminId int64) context.Context {
	ctx := context.Background()
	m := maps.Map{
		"timestamp": time.Now().Unix(),
		"type":      "admin",
		"userId":    adminId,
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
