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

// RPC客户端
type RPCClient struct {
	apiConfig *configs.APIConfig
	conns     []*grpc.ClientConn
}

// 构造新的RPC客户端
func NewRPCClient(apiConfig *configs.APIConfig) (*RPCClient, error) {
	if apiConfig == nil {
		return nil, errors.New("api config should not be nil")
	}

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

	return &RPCClient{
		apiConfig: apiConfig,
		conns:     conns,
	}, nil
}

func (this *RPCClient) AdminRPC() pb.AdminServiceClient {
	return pb.NewAdminServiceClient(this.pickConn())
}

func (this *RPCClient) NodeRPC() pb.NodeServiceClient {
	return pb.NewNodeServiceClient(this.pickConn())
}

func (this *RPCClient) NodeGrantRPC() pb.NodeGrantServiceClient {
	return pb.NewNodeGrantServiceClient(this.pickConn())
}

func (this *RPCClient) NodeClusterRPC() pb.NodeClusterServiceClient {
	return pb.NewNodeClusterServiceClient(this.pickConn())
}

func (this *RPCClient) NodeIPAddressRPC() pb.NodeIPAddressServiceClient {
	return pb.NewNodeIPAddressServiceClient(this.pickConn())
}

func (this *RPCClient) ServerRPC() pb.ServerServiceClient {
	return pb.NewServerServiceClient(this.pickConn())
}

func (this *RPCClient) APINodeRPC() pb.APINodeServiceClient {
	return pb.NewAPINodeServiceClient(this.pickConn())
}

func (this *RPCClient) OriginRPC() pb.OriginServiceClient {
	return pb.NewOriginServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPWebRPC() pb.HTTPWebServiceClient {
	return pb.NewHTTPWebServiceClient(this.pickConn())
}

func (this *RPCClient) ReverseProxyRPC() pb.ReverseProxyServiceClient {
	return pb.NewReverseProxyServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPGzipRPC() pb.HTTPGzipServiceClient {
	return pb.NewHTTPGzipServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPHeaderRPC() pb.HTTPHeaderServiceClient {
	return pb.NewHTTPHeaderServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPHeaderPolicyRPC() pb.HTTPHeaderPolicyServiceClient {
	return pb.NewHTTPHeaderPolicyServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPPageRPC() pb.HTTPPageServiceClient {
	return pb.NewHTTPPageServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPAccessLogPolicyRPC() pb.HTTPAccessLogPolicyServiceClient {
	return pb.NewHTTPAccessLogPolicyServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPCachePolicyRPC() pb.HTTPCachePolicyServiceClient {
	return pb.NewHTTPCachePolicyServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPFirewallPolicyRPC() pb.HTTPFirewallPolicyServiceClient {
	return pb.NewHTTPFirewallPolicyServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPLocationRPC() pb.HTTPLocationServiceClient {
	return pb.NewHTTPLocationServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPWebsocketRPC() pb.HTTPWebsocketServiceClient {
	return pb.NewHTTPWebsocketServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPRewriteRuleRPC() pb.HTTPRewriteRuleServiceClient {
	return pb.NewHTTPRewriteRuleServiceClient(this.pickConn())
}

func (this *RPCClient) SSLCertRPC() pb.SSLCertServiceClient {
	return pb.NewSSLCertServiceClient(this.pickConn())
}

func (this *RPCClient) SSLPolicyRPC() pb.SSLPolicyServiceClient {
	return pb.NewSSLPolicyServiceClient(this.pickConn())
}

func (this *RPCClient) SysSettingRPC() pb.SysSettingServiceClient {
	return pb.NewSysSettingServiceClient(this.pickConn())
}

// 构造上下文
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

// 随机选择一个连接
func (this *RPCClient) pickConn() *grpc.ClientConn {
	if len(this.conns) == 0 {
		return nil
	}
	return this.conns[rands.Int(0, len(this.conns)-1)]
}
