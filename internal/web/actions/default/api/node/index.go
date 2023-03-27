package node

import (
	"encoding/json"
	"fmt"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "index")
}

func (this *IndexAction) RunGet(params struct {
	NodeId int64
}) {
	nodeResp, err := this.RPC().APINodeRPC().FindEnabledAPINode(this.AdminContext(), &pb.FindEnabledAPINodeRequest{ApiNodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var node = nodeResp.ApiNode
	if node == nil {
		this.NotFound("apiNode", params.NodeId)
		return
	}

	// 监听地址
	var hasHTTPS = false
	var httpConfig = &serverconfigs.HTTPProtocolConfig{}
	if len(node.HttpJSON) > 0 {
		err = json.Unmarshal(node.HttpJSON, httpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	var httpsConfig = &serverconfigs.HTTPSProtocolConfig{}
	if len(node.HttpsJSON) > 0 {
		err = json.Unmarshal(node.HttpsJSON, httpsConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		hasHTTPS = len(httpsConfig.Listen) > 0
	}

	// 监听地址
	var listens = []*serverconfigs.NetworkAddressConfig{}
	listens = append(listens, httpConfig.Listen...)
	listens = append(listens, httpsConfig.Listen...)

	// 证书信息
	var certs = []*sslconfigs.SSLCertConfig{}
	if httpsConfig.SSLPolicyRef != nil && httpsConfig.SSLPolicyRef.SSLPolicyId > 0 {
		sslPolicyConfigResp, err := this.RPC().SSLPolicyRPC().FindEnabledSSLPolicyConfig(this.AdminContext(), &pb.FindEnabledSSLPolicyConfigRequest{
			SslPolicyId: httpsConfig.SSLPolicyRef.SSLPolicyId,
			IgnoreData:  true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var sslPolicyConfigJSON = sslPolicyConfigResp.SslPolicyJSON
		if len(sslPolicyConfigJSON) > 0 {
			var sslPolicy = &sslconfigs.SSLPolicy{}
			err = json.Unmarshal(sslPolicyConfigJSON, sslPolicy)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			certs = sslPolicy.Certs
		}
	}

	// 访问地址
	var accessAddrs = []*serverconfigs.NetworkAddressConfig{}
	if len(node.AccessAddrsJSON) > 0 {
		err = json.Unmarshal(node.AccessAddrsJSON, &accessAddrs)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// Rest地址
	var restAccessAddrs = []*serverconfigs.NetworkAddressConfig{}
	if node.RestIsOn {
		if len(node.RestHTTPJSON) > 0 {
			var httpConfig = &serverconfigs.HTTPProtocolConfig{}
			err = json.Unmarshal(node.RestHTTPJSON, httpConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if httpConfig.IsOn && len(httpConfig.Listen) > 0 {
				restAccessAddrs = append(restAccessAddrs, httpConfig.Listen...)
			}
		}

		if len(node.RestHTTPSJSON) > 0 {
			var httpsConfig = &serverconfigs.HTTPSProtocolConfig{}
			err = json.Unmarshal(node.RestHTTPSJSON, httpsConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if httpsConfig.IsOn && len(httpsConfig.Listen) > 0 {
				restAccessAddrs = append(restAccessAddrs, httpsConfig.Listen...)
			}

			if !hasHTTPS {
				hasHTTPS = len(httpsConfig.Listen) > 0
			}
		}
	}

	// 状态
	var status = &nodeconfigs.NodeStatus{}
	var statusIsValid = false
	this.Data["newVersion"] = ""
	if len(node.StatusJSON) > 0 {
		err = json.Unmarshal(node.StatusJSON, &status)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		if status.UpdatedAt >= time.Now().Unix()-300 {
			statusIsValid = true

			// 是否为新版本
			if stringutil.VersionCompare(status.BuildVersion, teaconst.APINodeVersion) < 0 {
				this.Data["newVersion"] = teaconst.APINodeVersion
			}
		}
	}

	this.Data["node"] = maps.Map{
		"id":              node.Id,
		"name":            node.Name,
		"description":     node.Description,
		"isOn":            node.IsOn,
		"listens":         listens,
		"accessAddrs":     accessAddrs,
		"restIsOn":        node.RestIsOn,
		"restAccessAddrs": restAccessAddrs,
		"hasHTTPS":        hasHTTPS,
		"certs":           certs,
		"isPrimary":       node.IsPrimary,
		"statusIsValid":   statusIsValid,
		"status": maps.Map{
			"isActive":             status.IsActive,
			"updatedAt":            status.UpdatedAt,
			"hostname":             status.Hostname,
			"cpuUsage":             status.CPUUsage,
			"cpuUsageText":         fmt.Sprintf("%.2f%%", status.CPUUsage*100),
			"memUsage":             status.MemoryUsage,
			"memUsageText":         fmt.Sprintf("%.2f%%", status.MemoryUsage*100),
			"connectionCount":      status.ConnectionCount,
			"buildVersion":         status.BuildVersion,
			"cpuPhysicalCount":     status.CPUPhysicalCount,
			"cpuLogicalCount":      status.CPULogicalCount,
			"load1m":               numberutils.FormatFloat2(status.Load1m),
			"load5m":               numberutils.FormatFloat2(status.Load5m),
			"load15m":              numberutils.FormatFloat2(status.Load15m),
			"cacheTotalDiskSize":   numberutils.FormatBytes(status.CacheTotalDiskSize),
			"cacheTotalMemorySize": numberutils.FormatBytes(status.CacheTotalMemorySize),
			"exePath":              status.ExePath,
		},
	}

	this.Show()
}
