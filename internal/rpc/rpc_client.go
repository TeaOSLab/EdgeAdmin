package rpc

import (
	"context"
	"encoding/base64"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/encrypt"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

type RPCClient struct {
	apiConfig                  *configs.APIConfig
	adminClients               []pb.AdminServiceClient
	nodeClients                []pb.NodeServiceClient
	nodeGrantClients           []pb.NodeGrantServiceClient
	nodeClusterClients         []pb.NodeClusterServiceClient
	nodeIPAddressClients       []pb.NodeIPAddressServiceClient
	serverClients              []pb.ServerServiceClient
	apiNodeClients             []pb.APINodeServiceClient
	originNodeClients          []pb.OriginServerServiceClient
	httpWebClients             []pb.HTTPWebServiceClient
	reverseProxyClients        []pb.ReverseProxyServiceClient
	httpGzipClients            []pb.HTTPGzipServiceClient
	httpHeaderPolicyClients    []pb.HTTPHeaderPolicyServiceClient
	httpHeaderClients          []pb.HTTPHeaderServiceClient
	httpPageClients            []pb.HTTPPageServiceClient
	httpAccessLogPolicyClients []pb.HTTPAccessLogPolicyServiceClient
	httpCachePolicyClients     []pb.HTTPCachePolicyServiceClient
	httpFirewallPolicyClients  []pb.HTTPFirewallPolicyServiceClient
	httpLocationClients        []pb.HTTPLocationServiceClient
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
	originNodeClients := []pb.OriginServerServiceClient{}
	httpWebClients := []pb.HTTPWebServiceClient{}
	reverseProxyClients := []pb.ReverseProxyServiceClient{}
	httpGzipClients := []pb.HTTPGzipServiceClient{}
	httpHeaderPolicyClients := []pb.HTTPHeaderPolicyServiceClient{}
	httpHeaderClients := []pb.HTTPHeaderServiceClient{}
	httpPageClients := []pb.HTTPPageServiceClient{}
	httpAccessLogPolicyClients := []pb.HTTPAccessLogPolicyServiceClient{}
	httpCachePolicyClients := []pb.HTTPCachePolicyServiceClient{}
	httpFirewallPolicyClients := []pb.HTTPFirewallPolicyServiceClient{}
	httpLocationClients := []pb.HTTPLocationServiceClient{}

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
		originNodeClients = append(originNodeClients, pb.NewOriginServerServiceClient(conn))
		httpWebClients = append(httpWebClients, pb.NewHTTPWebServiceClient(conn))
		reverseProxyClients = append(reverseProxyClients, pb.NewReverseProxyServiceClient(conn))
		httpGzipClients = append(httpGzipClients, pb.NewHTTPGzipServiceClient(conn))
		httpHeaderPolicyClients = append(httpHeaderPolicyClients, pb.NewHTTPHeaderPolicyServiceClient(conn))
		httpHeaderClients = append(httpHeaderClients, pb.NewHTTPHeaderServiceClient(conn))
		httpPageClients = append(httpPageClients, pb.NewHTTPPageServiceClient(conn))
		httpAccessLogPolicyClients = append(httpAccessLogPolicyClients, pb.NewHTTPAccessLogPolicyServiceClient(conn))
		httpCachePolicyClients = append(httpCachePolicyClients, pb.NewHTTPCachePolicyServiceClient(conn))
		httpFirewallPolicyClients = append(httpFirewallPolicyClients, pb.NewHTTPFirewallPolicyServiceClient(conn))
		httpLocationClients = append(httpLocationClients, pb.NewHTTPLocationServiceClient(conn))
	}

	return &RPCClient{
		apiConfig:                  apiConfig,
		adminClients:               adminClients,
		nodeClients:                nodeClients,
		nodeGrantClients:           nodeGrantClients,
		nodeClusterClients:         nodeClusterClients,
		nodeIPAddressClients:       nodeIPAddressClients,
		serverClients:              serverClients,
		apiNodeClients:             apiNodeClients,
		originNodeClients:          originNodeClients,
		httpWebClients:             httpWebClients,
		reverseProxyClients:        reverseProxyClients,
		httpGzipClients:            httpGzipClients,
		httpHeaderPolicyClients:    httpHeaderPolicyClients,
		httpHeaderClients:          httpHeaderClients,
		httpPageClients:            httpPageClients,
		httpAccessLogPolicyClients: httpAccessLogPolicyClients,
		httpCachePolicyClients:     httpCachePolicyClients,
		httpFirewallPolicyClients:  httpFirewallPolicyClients,
		httpLocationClients:        httpLocationClients,
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
	if len(this.apiNodeClients) > 0 {
		return this.apiNodeClients[rands.Int(0, len(this.apiNodeClients)-1)]
	}
	return nil
}

func (this *RPCClient) OriginServerRPC() pb.OriginServerServiceClient {
	if len(this.originNodeClients) > 0 {
		return this.originNodeClients[rands.Int(0, len(this.originNodeClients)-1)]
	}
	return nil
}

func (this *RPCClient) HTTPWebRPC() pb.HTTPWebServiceClient {
	if len(this.httpWebClients) > 0 {
		return this.httpWebClients[rands.Int(0, len(this.httpWebClients)-1)]
	}
	return nil
}

func (this *RPCClient) ReverseProxyRPC() pb.ReverseProxyServiceClient {
	if len(this.reverseProxyClients) > 0 {
		return this.reverseProxyClients[rands.Int(0, len(this.reverseProxyClients)-1)]
	}
	return nil
}

func (this *RPCClient) HTTPGzipRPC() pb.HTTPGzipServiceClient {
	if len(this.httpGzipClients) > 0 {
		return this.httpGzipClients[rands.Int(0, len(this.httpGzipClients)-1)]
	}
	return nil
}

func (this *RPCClient) HTTPHeaderRPC() pb.HTTPHeaderServiceClient {
	if len(this.httpHeaderClients) > 0 {
		return this.httpHeaderClients[rands.Int(0, len(this.httpHeaderClients)-1)]
	}
	return nil
}

func (this *RPCClient) HTTPHeaderPolicyRPC() pb.HTTPHeaderPolicyServiceClient {
	if len(this.httpHeaderPolicyClients) > 0 {
		return this.httpHeaderPolicyClients[rands.Int(0, len(this.httpHeaderPolicyClients)-1)]
	}
	return nil
}

func (this *RPCClient) HTTPPageRPC() pb.HTTPPageServiceClient {
	if len(this.httpPageClients) > 0 {
		return this.httpPageClients[rands.Int(0, len(this.httpPageClients)-1)]
	}
	return nil
}

func (this *RPCClient) HTTPAccessLogPolicyRPC() pb.HTTPAccessLogPolicyServiceClient {
	if len(this.httpAccessLogPolicyClients) > 0 {
		return this.httpAccessLogPolicyClients[rands.Int(0, len(this.httpAccessLogPolicyClients)-1)]
	}
	return nil
}

func (this *RPCClient) HTTPCachePolicyRPC() pb.HTTPCachePolicyServiceClient {
	if len(this.httpCachePolicyClients) > 0 {
		return this.httpCachePolicyClients[rands.Int(0, len(this.httpCachePolicyClients)-1)]
	}
	return nil
}

func (this *RPCClient) HTTPFirewallPolicyRPC() pb.HTTPFirewallPolicyServiceClient {
	if len(this.httpFirewallPolicyClients) > 0 {
		return this.httpFirewallPolicyClients[rands.Int(0, len(this.httpFirewallPolicyClients)-1)]
	}
	return nil
}

func (this *RPCClient) HTTPLocationRPC() pb.HTTPLocationServiceClient {
	if len(this.httpLocationClients) > 0 {
		return this.httpLocationClients[rands.Int(0, len(this.httpLocationClients)-1)]
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
