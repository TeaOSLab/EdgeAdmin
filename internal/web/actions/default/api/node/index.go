package node

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/maps"
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
	nodeResp, err := this.RPC().APINodeRPC().FindEnabledAPINode(this.AdminContext(), &pb.FindEnabledAPINodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	node := nodeResp.Node
	if node == nil {
		this.NotFound("apiNode", params.NodeId)
		return
	}

	// 监听地址
	var hasHTTPS = false
	httpConfig := &serverconfigs.HTTPProtocolConfig{}
	if len(node.HttpJSON) > 0 {
		err = json.Unmarshal(node.HttpJSON, httpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	httpsConfig := &serverconfigs.HTTPSProtocolConfig{}
	if len(node.HttpsJSON) > 0 {
		err = json.Unmarshal(node.HttpsJSON, httpsConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		hasHTTPS = len(httpsConfig.Listen) > 0
	}

	// 监听地址
	listens := []*serverconfigs.NetworkAddressConfig{}
	listens = append(listens, httpConfig.Listen...)
	listens = append(listens, httpsConfig.Listen...)

	// 证书信息
	certs := []*sslconfigs.SSLCertConfig{}
	if httpsConfig.SSLPolicyRef != nil && httpsConfig.SSLPolicyRef.SSLPolicyId > 0 {
		sslPolicyConfigResp, err := this.RPC().SSLPolicyRPC().FindEnabledSSLPolicyConfig(this.AdminContext(), &pb.FindEnabledSSLPolicyConfigRequest{SslPolicyId: httpsConfig.SSLPolicyRef.SSLPolicyId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		sslPolicyConfigJSON := sslPolicyConfigResp.SslPolicyJSON
		if len(sslPolicyConfigJSON) > 0 {
			sslPolicy := &sslconfigs.SSLPolicy{}
			err = json.Unmarshal(sslPolicyConfigJSON, sslPolicy)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			certs = sslPolicy.Certs
		}
	}

	// 访问地址
	accessAddrs := []*serverconfigs.NetworkAddressConfig{}
	if len(node.AccessAddrsJSON) > 0 {
		err = json.Unmarshal(node.AccessAddrsJSON, &accessAddrs)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// Rest地址
	restAccessAddrs := []*serverconfigs.NetworkAddressConfig{}
	if node.RestIsOn {
		if len(node.RestHTTPJSON) > 0 {
			httpConfig := &serverconfigs.HTTPProtocolConfig{}
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
			httpsConfig := &serverconfigs.HTTPSProtocolConfig{}
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
	}

	this.Show()
}
