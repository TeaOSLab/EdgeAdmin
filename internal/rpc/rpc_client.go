package rpc

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/encrypt"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"net/url"
	"sync"
	"time"
)

// RPCClient RPC客户端
type RPCClient struct {
	apiConfig *configs.APIConfig
	conns     []*grpc.ClientConn

	locker sync.Mutex
}

// NewRPCClient 构造新的RPC客户端
func NewRPCClient(apiConfig *configs.APIConfig) (*RPCClient, error) {
	if apiConfig == nil {
		return nil, errors.New("api config should not be nil")
	}

	client := &RPCClient{
		apiConfig: apiConfig,
	}

	err := client.init()
	if err != nil {
		return nil, err
	}

	// 设置RPC
	dao.SetRPC(client)

	return client, nil
}

func (this *RPCClient) APITokenRPC() pb.APITokenServiceClient {
	return pb.NewAPITokenServiceClient(this.pickConn())
}

func (this *RPCClient) AdminRPC() pb.AdminServiceClient {
	return pb.NewAdminServiceClient(this.pickConn())
}

func (this *RPCClient) NodeRPC() pb.NodeServiceClient {
	return pb.NewNodeServiceClient(this.pickConn())
}

func (this *RPCClient) NodeLogRPC() pb.NodeLogServiceClient {
	return pb.NewNodeLogServiceClient(this.pickConn())
}

func (this *RPCClient) NodeGrantRPC() pb.NodeGrantServiceClient {
	return pb.NewNodeGrantServiceClient(this.pickConn())
}

func (this *RPCClient) NodeClusterRPC() pb.NodeClusterServiceClient {
	return pb.NewNodeClusterServiceClient(this.pickConn())
}

func (this *RPCClient) NodeClusterFirewallActionRPC() pb.NodeClusterFirewallActionServiceClient {
	return pb.NewNodeClusterFirewallActionServiceClient(this.pickConn())
}

func (this *RPCClient) NodeGroupRPC() pb.NodeGroupServiceClient {
	return pb.NewNodeGroupServiceClient(this.pickConn())
}

func (this *RPCClient) NodeRegionRPC() pb.NodeRegionServiceClient {
	return pb.NewNodeRegionServiceClient(this.pickConn())
}

func (this *RPCClient) NodePriceItemRPC() pb.NodePriceItemServiceClient {
	return pb.NewNodePriceItemServiceClient(this.pickConn())
}

func (this *RPCClient) NodeIPAddressRPC() pb.NodeIPAddressServiceClient {
	return pb.NewNodeIPAddressServiceClient(this.pickConn())
}

func (this *RPCClient) NodeValueRPC() pb.NodeValueServiceClient {
	return pb.NewNodeValueServiceClient(this.pickConn())
}

func (this *RPCClient) NodeThresholdRPC() pb.NodeThresholdServiceClient {
	return pb.NewNodeThresholdServiceClient(this.pickConn())
}

func (this *RPCClient) ServerRPC() pb.ServerServiceClient {
	return pb.NewServerServiceClient(this.pickConn())
}

func (this *RPCClient) ServerClientSystemMonthlyStatRPC() pb.ServerClientSystemMonthlyStatServiceClient {
	return pb.NewServerClientSystemMonthlyStatServiceClient(this.pickConn())
}

func (this *RPCClient) ServerClientBrowserMonthlyStatRPC() pb.ServerClientBrowserMonthlyStatServiceClient {
	return pb.NewServerClientBrowserMonthlyStatServiceClient(this.pickConn())
}

func (this *RPCClient) ServerRegionCountryMonthlyStatRPC() pb.ServerRegionCountryMonthlyStatServiceClient {
	return pb.NewServerRegionCountryMonthlyStatServiceClient(this.pickConn())
}

func (this *RPCClient) ServerRegionProvinceMonthlyStatRPC() pb.ServerRegionProvinceMonthlyStatServiceClient {
	return pb.NewServerRegionProvinceMonthlyStatServiceClient(this.pickConn())
}

func (this *RPCClient) ServerRegionCityMonthlyStatRPC() pb.ServerRegionCityMonthlyStatServiceClient {
	return pb.NewServerRegionCityMonthlyStatServiceClient(this.pickConn())
}

func (this *RPCClient) ServerRegionProviderMonthlyStatRPC() pb.ServerRegionProviderMonthlyStatServiceClient {
	return pb.NewServerRegionProviderMonthlyStatServiceClient(this.pickConn())
}

func (this *RPCClient) ServerHTTPFirewallDailyStatRPC() pb.ServerHTTPFirewallDailyStatServiceClient {
	return pb.NewServerHTTPFirewallDailyStatServiceClient(this.pickConn())
}

func (this *RPCClient) ServerDailyStatRPC() pb.ServerDailyStatServiceClient {
	return pb.NewServerDailyStatServiceClient(this.pickConn())
}

func (this *RPCClient) ServerGroupRPC() pb.ServerGroupServiceClient {
	return pb.NewServerGroupServiceClient(this.pickConn())
}

func (this *RPCClient) APINodeRPC() pb.APINodeServiceClient {
	return pb.NewAPINodeServiceClient(this.pickConn())
}

func (this *RPCClient) UserNodeRPC() pb.UserNodeServiceClient {
	return pb.NewUserNodeServiceClient(this.pickConn())
}

func (this *RPCClient) DBNodeRPC() pb.DBNodeServiceClient {
	return pb.NewDBNodeServiceClient(this.pickConn())
}

func (this *RPCClient) MonitorNodeRPC() pb.MonitorNodeServiceClient {
	return pb.NewMonitorNodeServiceClient(this.pickConn())
}

func (this *RPCClient) DBRPC() pb.DBServiceClient {
	return pb.NewDBServiceClient(this.pickConn())
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

func (this *RPCClient) HTTPFirewallRuleGroupRPC() pb.HTTPFirewallRuleGroupServiceClient {
	return pb.NewHTTPFirewallRuleGroupServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPFirewallRuleSetRPC() pb.HTTPFirewallRuleSetServiceClient {
	return pb.NewHTTPFirewallRuleSetServiceClient(this.pickConn())
}

func (this *RPCClient) FirewallRPC() pb.FirewallServiceClient {
	return pb.NewFirewallServiceClient(this.pickConn())
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

// HTTPAccessLogRPC 访问日志
func (this *RPCClient) HTTPAccessLogRPC() pb.HTTPAccessLogServiceClient {
	return pb.NewHTTPAccessLogServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPFastcgiRPC() pb.HTTPFastcgiServiceClient {
	return pb.NewHTTPFastcgiServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPAuthPolicyRPC() pb.HTTPAuthPolicyServiceClient {
	return pb.NewHTTPAuthPolicyServiceClient(this.pickConn())
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

func (this *RPCClient) MessageRPC() pb.MessageServiceClient {
	return pb.NewMessageServiceClient(this.pickConn())
}

func (this *RPCClient) MessageRecipientGroupRPC() pb.MessageRecipientGroupServiceClient {
	return pb.NewMessageRecipientGroupServiceClient(this.pickConn())
}

func (this *RPCClient) MessageRecipientRPC() pb.MessageRecipientServiceClient {
	return pb.NewMessageRecipientServiceClient(this.pickConn())
}

func (this *RPCClient) MessageMediaRPC() pb.MessageMediaServiceClient {
	return pb.NewMessageMediaServiceClient(this.pickConn())
}

func (this *RPCClient) MessageMediaInstanceRPC() pb.MessageMediaInstanceServiceClient {
	return pb.NewMessageMediaInstanceServiceClient(this.pickConn())
}

func (this *RPCClient) MessageTaskRPC() pb.MessageTaskServiceClient {
	return pb.NewMessageTaskServiceClient(this.pickConn())
}

func (this *RPCClient) MessageTaskLogRPC() pb.MessageTaskLogServiceClient {
	return pb.NewMessageTaskLogServiceClient(this.pickConn())
}

func (this *RPCClient) MessageReceiverRPC() pb.MessageReceiverServiceClient {
	return pb.NewMessageReceiverServiceClient(this.pickConn())
}

func (this *RPCClient) IPLibraryRPC() pb.IPLibraryServiceClient {
	return pb.NewIPLibraryServiceClient(this.pickConn())
}

func (this *RPCClient) IPListRPC() pb.IPListServiceClient {
	return pb.NewIPListServiceClient(this.pickConn())
}

func (this *RPCClient) IPItemRPC() pb.IPItemServiceClient {
	return pb.NewIPItemServiceClient(this.pickConn())
}

func (this *RPCClient) FileRPC() pb.FileServiceClient {
	return pb.NewFileServiceClient(this.pickConn())
}

func (this *RPCClient) FileChunkRPC() pb.FileChunkServiceClient {
	return pb.NewFileChunkServiceClient(this.pickConn())
}

func (this *RPCClient) RegionCountryRPC() pb.RegionCountryServiceClient {
	return pb.NewRegionCountryServiceClient(this.pickConn())
}

func (this *RPCClient) RegionProvinceRPC() pb.RegionProvinceServiceClient {
	return pb.NewRegionProvinceServiceClient(this.pickConn())
}

func (this *RPCClient) LogRPC() pb.LogServiceClient {
	return pb.NewLogServiceClient(this.pickConn())
}

func (this *RPCClient) DNSProviderRPC() pb.DNSProviderServiceClient {
	return pb.NewDNSProviderServiceClient(this.pickConn())
}

func (this *RPCClient) DNSDomainRPC() pb.DNSDomainServiceClient {
	return pb.NewDNSDomainServiceClient(this.pickConn())
}

func (this *RPCClient) DNSRPC() pb.DNSServiceClient {
	return pb.NewDNSServiceClient(this.pickConn())
}

func (this *RPCClient) DNSTaskRPC() pb.DNSTaskServiceClient {
	return pb.NewDNSTaskServiceClient(this.pickConn())
}

func (this *RPCClient) ACMEUserRPC() pb.ACMEUserServiceClient {
	return pb.NewACMEUserServiceClient(this.pickConn())
}

func (this *RPCClient) ACMETaskRPC() pb.ACMETaskServiceClient {
	return pb.NewACMETaskServiceClient(this.pickConn())
}

func (this *RPCClient) UserRPC() pb.UserServiceClient {
	return pb.NewUserServiceClient(this.pickConn())
}

func (this *RPCClient) UserBillRPC() pb.UserBillServiceClient {
	return pb.NewUserBillServiceClient(this.pickConn())
}

func (this *RPCClient) UserAccessKeyRPC() pb.UserAccessKeyServiceClient {
	return pb.NewUserAccessKeyServiceClient(this.pickConn())
}

func (this *RPCClient) LoginRPC() pb.LoginServiceClient {
	return pb.NewLoginServiceClient(this.pickConn())
}

func (this *RPCClient) NodeTaskRPC() pb.NodeTaskServiceClient {
	return pb.NewNodeTaskServiceClient(this.pickConn())
}

func (this *RPCClient) AuthorityKeyRPC() pb.AuthorityKeyServiceClient {
	return pb.NewAuthorityKeyServiceClient(this.pickConn())
}

func (this *RPCClient) AuthorityNodeRPC() pb.AuthorityNodeServiceClient {
	return pb.NewAuthorityNodeServiceClient(this.pickConn())
}

func (this *RPCClient) LatestItemRPC() pb.LatestItemServiceClient {
	return pb.NewLatestItemServiceClient(this.pickConn())
}

func (this *RPCClient) NSClusterRPC() pb.NSClusterServiceClient {
	return pb.NewNSClusterServiceClient(this.pickConn())
}

func (this *RPCClient) NSNodeRPC() pb.NSNodeServiceClient {
	return pb.NewNSNodeServiceClient(this.pickConn())
}

func (this *RPCClient) NSDomainRPC() pb.NSDomainServiceClient {
	return pb.NewNSDomainServiceClient(this.pickConn())
}

func (this *RPCClient) NSRecordRPC() pb.NSRecordServiceClient {
	return pb.NewNSRecordServiceClient(this.pickConn())
}

func (this *RPCClient) NSKeyRPC() pb.NSKeyServiceClient {
	return pb.NewNSKeyServiceClient(this.pickConn())
}

func (this *RPCClient) NSRouteRPC() pb.NSRouteServiceClient {
	return pb.NewNSRouteServiceClient(this.pickConn())
}

func (this *RPCClient) NSAccessLogRPC() pb.NSAccessLogServiceClient {
	return pb.NewNSAccessLogServiceClient(this.pickConn())
}

func (this *RPCClient) NSRPC() pb.NSServiceClient {
	return pb.NewNSServiceClient(this.pickConn())
}

func (this *RPCClient) MetricItemRPC() pb.MetricItemServiceClient {
	return pb.NewMetricItemServiceClient(this.pickConn())
}

func (this *RPCClient) MetricStatRPC() pb.MetricStatServiceClient {
	return pb.NewMetricStatServiceClient(this.pickConn())
}

func (this *RPCClient) MetricChartRPC() pb.MetricChartServiceClient {
	return pb.NewMetricChartServiceClient(this.pickConn())
}

func (this *RPCClient) NodeClusterMetricItemRPC() pb.NodeClusterMetricItemServiceClient {
	return pb.NewNodeClusterMetricItemServiceClient(this.pickConn())
}

func (this *RPCClient) ServerStatBoardRPC() pb.ServerStatBoardServiceClient {
	return pb.NewServerStatBoardServiceClient(this.pickConn())
}

func (this *RPCClient) ServerStatBoardChartRPC() pb.ServerStatBoardChartServiceClient {
	return pb.NewServerStatBoardChartServiceClient(this.pickConn())
}

// Context 构造Admin上下文
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

// APIContext 构造API上下文
func (this *RPCClient) APIContext(apiNodeId int64) context.Context {
	ctx := context.Background()
	m := maps.Map{
		"timestamp": time.Now().Unix(),
		"type":      "api",
		"userId":    apiNodeId,
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

// UpdateConfig 修改配置
func (this *RPCClient) UpdateConfig(config *configs.APIConfig) error {
	this.apiConfig = config
	return this.init()
}

// 初始化
func (this *RPCClient) init() error {
	// 重新连接
	conns := []*grpc.ClientConn{}
	for _, endpoint := range this.apiConfig.RPC.Endpoints {
		u, err := url.Parse(endpoint)
		if err != nil {
			return errors.New("parse endpoint failed: " + err.Error())
		}
		var conn *grpc.ClientConn
		if u.Scheme == "http" {
			conn, err = grpc.Dial(u.Host, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(128*1024*1024)))
		} else if u.Scheme == "https" {
			conn, err = grpc.Dial(u.Host, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(128*1024*1024)))
		} else {
			return errors.New("parse endpoint failed: invalid scheme '" + u.Scheme + "'")
		}
		if err != nil {
			return err
		}
		conns = append(conns, conn)
	}
	if len(conns) == 0 {
		return errors.New("[RPC]no available endpoints")
	}

	// 这里不需要加锁，因为会和pickConn冲突
	this.conns = conns
	return nil
}

// 随机选择一个连接
func (this *RPCClient) pickConn() *grpc.ClientConn {
	this.locker.Lock()
	defer this.locker.Unlock()

	// 检查连接状态
	if len(this.conns) > 0 {
		availableConns := []*grpc.ClientConn{}
		for _, state := range []connectivity.State{connectivity.Ready, connectivity.Idle, connectivity.Connecting} {
			for _, conn := range this.conns {
				if conn.GetState() == state {
					availableConns = append(availableConns, conn)
				}
			}
			if len(availableConns) > 0 {
				break
			}
		}

		if len(availableConns) > 0 {
			return availableConns[rands.Int(0, len(availableConns)-1)]
		}

		// 关闭
		for _, conn := range this.conns {
			_ = conn.Close()
		}
	}

	// 重新初始化
	err := this.init()
	if err != nil {
		// 错误提示已经在构造对象时打印过，所以这里不再重复打印
		return nil
	}

	if len(this.conns) == 0 {
		return nil
	}

	return this.conns[rands.Int(0, len(this.conns)-1)]
}
