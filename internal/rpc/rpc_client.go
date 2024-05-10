package rpc

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/encrypt"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/encoding/gzip"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"
)

// RPCClient RPC客户端
type RPCClient struct {
	apiConfig *configs.APIConfig
	conns     []*grpc.ClientConn

	locker sync.RWMutex
}

// NewRPCClient 构造新的RPC客户端
func NewRPCClient(apiConfig *configs.APIConfig, isPrimary bool) (*RPCClient, error) {
	if apiConfig == nil {
		return nil, errors.New("api config should not be nil")
	}

	var client = &RPCClient{
		apiConfig: apiConfig,
	}

	err := client.init()
	if err != nil {
		return nil, err
	}

	// 设置RPC
	if isPrimary {
		dao.SetRPC(client)
	}

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

func (this *RPCClient) NodeLoginRPC() pb.NodeLoginServiceClient {
	return pb.NewNodeLoginServiceClient(this.pickConn())
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

func (this *RPCClient) NodeIPAddressRPC() pb.NodeIPAddressServiceClient {
	return pb.NewNodeIPAddressServiceClient(this.pickConn())
}

func (this *RPCClient) NodeIPAddressLogRPC() pb.NodeIPAddressLogServiceClient {
	return pb.NewNodeIPAddressLogServiceClient(this.pickConn())
}

func (this *RPCClient) NodeIPAddressThresholdRPC() pb.NodeIPAddressThresholdServiceClient {
	return pb.NewNodeIPAddressThresholdServiceClient(this.pickConn())
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

func (this *RPCClient) ServerBandwidthStatRPC() pb.ServerBandwidthStatServiceClient {
	return pb.NewServerBandwidthStatServiceClient(this.pickConn())
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

func (this *RPCClient) APIMethodStatRPC() pb.APIMethodStatServiceClient {
	return pb.NewAPIMethodStatServiceClient(this.pickConn())
}

func (this *RPCClient) DBNodeRPC() pb.DBNodeServiceClient {
	return pb.NewDBNodeServiceClient(this.pickConn())
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

func (this *RPCClient) HTTPCacheTaskRPC() pb.HTTPCacheTaskServiceClient {
	return pb.NewHTTPCacheTaskServiceClient(this.pickConn())
}

func (this *RPCClient) HTTPCacheTaskKeyRPC() pb.HTTPCacheTaskKeyServiceClient {
	return pb.NewHTTPCacheTaskKeyServiceClient(this.pickConn())
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

func (this *RPCClient) IPLibraryRPC() pb.IPLibraryServiceClient {
	return pb.NewIPLibraryServiceClient(this.pickConn())
}

func (this *RPCClient) IPLibraryFileRPC() pb.IPLibraryFileServiceClient {
	return pb.NewIPLibraryFileServiceClient(this.pickConn())
}

func (this *RPCClient) IPLibraryArtifactRPC() pb.IPLibraryArtifactServiceClient {
	return pb.NewIPLibraryArtifactServiceClient(this.pickConn())
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

func (this *RPCClient) RegionCityRPC() pb.RegionCityServiceClient {
	return pb.NewRegionCityServiceClient(this.pickConn())
}

func (this *RPCClient) RegionTownRPC() pb.RegionTownServiceClient {
	return pb.NewRegionTownServiceClient(this.pickConn())
}

func (this *RPCClient) RegionProviderRPC() pb.RegionProviderServiceClient {
	return pb.NewRegionProviderServiceClient(this.pickConn())
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

func (this *RPCClient) ACMEProviderRPC() pb.ACMEProviderServiceClient {
	return pb.NewACMEProviderServiceClient(this.pickConn())
}

func (this *RPCClient) ACMEProviderAccountRPC() pb.ACMEProviderAccountServiceClient {
	return pb.NewACMEProviderAccountServiceClient(this.pickConn())
}

func (this *RPCClient) UserRPC() pb.UserServiceClient {
	return pb.NewUserServiceClient(this.pickConn())
}

func (this *RPCClient) UserAccessKeyRPC() pb.UserAccessKeyServiceClient {
	return pb.NewUserAccessKeyServiceClient(this.pickConn())
}

func (this *RPCClient) UserIdentityRPC() pb.UserIdentityServiceClient {
	return pb.NewUserIdentityServiceClient(this.pickConn())
}

func (this *RPCClient) LoginRPC() pb.LoginServiceClient {
	return pb.NewLoginServiceClient(this.pickConn())
}

func (this *RPCClient) LoginSessionRPC() pb.LoginSessionServiceClient {
	return pb.NewLoginSessionServiceClient(this.pickConn())
}

func (this *RPCClient) LoginTicketRPC() pb.LoginTicketServiceClient {
	return pb.NewLoginTicketServiceClient(this.pickConn())
}

func (this *RPCClient) NodeTaskRPC() pb.NodeTaskServiceClient {
	return pb.NewNodeTaskServiceClient(this.pickConn())
}

func (this *RPCClient) LatestItemRPC() pb.LatestItemServiceClient {
	return pb.NewLatestItemServiceClient(this.pickConn())
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

func (this *RPCClient) ServerDomainHourlyStatRPC() pb.ServerDomainHourlyStatServiceClient {
	return pb.NewServerDomainHourlyStatServiceClient(this.pickConn())
}

func (this *RPCClient) ServerStatBoardChartRPC() pb.ServerStatBoardChartServiceClient {
	return pb.NewServerStatBoardChartServiceClient(this.pickConn())
}

func (this *RPCClient) TrafficDailyStatRPC() pb.TrafficDailyStatServiceClient {
	return pb.NewTrafficDailyStatServiceClient(this.pickConn())
}

// Context 构造Admin上下文
func (this *RPCClient) Context(adminId int64) context.Context {
	var ctx = context.Background()
	var m = maps.Map{
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
	var token = base64.StdEncoding.EncodeToString(data)
	ctx = metadata.AppendToOutgoingContext(ctx, "nodeId", this.apiConfig.NodeId, "token", token)
	return ctx
}

// APIContext 构造API上下文
func (this *RPCClient) APIContext(apiNodeId int64) context.Context {
	var ctx = context.Background()
	var m = maps.Map{
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
	var token = base64.StdEncoding.EncodeToString(data)
	ctx = metadata.AppendToOutgoingContext(ctx, "nodeId", this.apiConfig.NodeId, "token", token)
	return ctx
}

// UpdateConfig 修改配置
func (this *RPCClient) UpdateConfig(config *configs.APIConfig) error {
	this.apiConfig = config

	this.locker.Lock()
	err := this.init()
	this.locker.Unlock()
	return err
}

// 初始化
func (this *RPCClient) init() error {
	// 当前的IP地址
	var localIPAddrs = this.localIPAddrs()

	// 重新连接
	var conns = []*grpc.ClientConn{}
	for _, endpoint := range this.apiConfig.RPCEndpoints {
		u, err := url.Parse(endpoint)
		if err != nil {
			return fmt.Errorf("parse endpoint failed: %w", err)
		}

		var apiHost = u.Host

		// 如果本机，则将地址修改为回路地址
		if lists.ContainsString(localIPAddrs, u.Hostname()) {
			if strings.Contains(apiHost, "[") { // IPv6 [host]:port
				apiHost = "[::1]"
			} else {
				apiHost = "127.0.0.1"
			}
			var port = u.Port()
			if len(port) > 0 {
				apiHost += ":" + port
			}
		}

		var conn *grpc.ClientConn
		var callOptions = grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(128<<20),
			grpc.MaxCallSendMsgSize(128<<20),
			grpc.UseCompressor(gzip.Name),
		)
		var keepaliveParams = grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 30 * time.Second,
		})
		if u.Scheme == "http" {
			conn, err = grpc.Dial(apiHost, grpc.WithTransportCredentials(insecure.NewCredentials()), callOptions, keepaliveParams)
		} else if u.Scheme == "https" {
			conn, err = grpc.Dial(apiHost, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				InsecureSkipVerify: true,
			})), callOptions, keepaliveParams)
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
	var countConns = len(this.conns)
	if countConns > 0 {
		if countConns == 1 {
			return this.conns[0]
		}
		for _, state := range []connectivity.State{
			connectivity.Ready,
			connectivity.Idle,
			connectivity.Connecting,
			connectivity.TransientFailure,
		} {
			var availableConns = []*grpc.ClientConn{}
			for _, conn := range this.conns {
				if conn.GetState() == state {
					availableConns = append(availableConns, conn)
				}
			}
			if len(availableConns) > 0 {
				return this.randConn(availableConns)
			}
		}
	}

	return this.randConn(this.conns)
}

// Close 关闭
func (this *RPCClient) Close() error {
	this.locker.Lock()
	defer this.locker.Unlock()

	var lastErr error
	for _, conn := range this.conns {
		var err = conn.Close()
		if err != nil {
			lastErr = err
			continue
		}
	}

	return lastErr
}

func (this *RPCClient) localIPAddrs() []string {
	localInterfaceAddrs, err := net.InterfaceAddrs()
	var localIPAddrs = []string{}
	if err == nil {
		for _, addr := range localInterfaceAddrs {
			var addrString = addr.String()
			var index = strings.Index(addrString, "/")
			if index > 0 {
				localIPAddrs = append(localIPAddrs, addrString[:index])
			}
		}
	}
	return localIPAddrs
}

func (this *RPCClient) randConn(conns []*grpc.ClientConn) *grpc.ClientConn {
	var l = len(conns)
	if l == 0 {
		return nil
	}
	if l == 1 {
		return conns[0]
	}
	return conns[rands.Int(0, l-1)]
}
