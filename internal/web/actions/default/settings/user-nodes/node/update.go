package node

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	NodeId int64
}) {
	nodeResp, err := this.RPC().UserNodeRPC().FindEnabledUserNode(this.AdminContext(), &pb.FindEnabledUserNodeRequest{
		NodeId: params.NodeId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	node := nodeResp.Node
	if node == nil {
		this.WriteString("要操作的节点不存在")
		return
	}

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
	}

	// 监听地址
	listens := []*serverconfigs.NetworkAddressConfig{}
	listens = append(listens, httpConfig.Listen...)
	listens = append(listens, httpsConfig.Listen...)

	// 证书信息
	certs := []*sslconfigs.SSLCertConfig{}
	sslPolicyId := int64(0)
	if httpsConfig.SSLPolicyRef != nil && httpsConfig.SSLPolicyRef.SSLPolicyId > 0 {
		sslPolicyConfigResp, err := this.RPC().SSLPolicyRPC().FindEnabledSSLPolicyConfig(this.AdminContext(), &pb.FindEnabledSSLPolicyConfigRequest{SslPolicyId: httpsConfig.SSLPolicyRef.SSLPolicyId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		sslPolicyConfigJSON := sslPolicyConfigResp.SslPolicyJSON
		if len(sslPolicyConfigJSON) > 0 {
			sslPolicyId = httpsConfig.SSLPolicyRef.SSLPolicyId

			sslPolicy := &sslconfigs.SSLPolicy{}
			err = json.Unmarshal(sslPolicyConfigJSON, sslPolicy)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			certs = sslPolicy.Certs
		}
	}

	accessAddrs := []*serverconfigs.NetworkAddressConfig{}
	if len(node.AccessAddrsJSON) > 0 {
		err = json.Unmarshal(node.AccessAddrsJSON, &accessAddrs)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["node"] = maps.Map{
		"id":          node.Id,
		"name":        node.Name,
		"description": node.Description,
		"isOn":        node.IsOn,
		"listens":     listens,
		"certs":       certs,
		"sslPolicyId": sslPolicyId,
		"accessAddrs": accessAddrs,
	}

	this.Show()
}

// 保存基础设置
func (this *UpdateAction) RunPost(params struct {
	NodeId          int64
	Name            string
	SslPolicyId     int64
	ListensJSON     []byte
	CertIdsJSON     []byte
	AccessAddrsJSON []byte
	Description     string
	IsOn            bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入用户节点名称")

	httpConfig := &serverconfigs.HTTPProtocolConfig{}
	httpsConfig := &serverconfigs.HTTPSProtocolConfig{}

	// 监听地址
	listens := []*serverconfigs.NetworkAddressConfig{}
	err := json.Unmarshal(params.ListensJSON, &listens)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if len(listens) == 0 {
		this.Fail("请添加至少一个进程监听地址")
	}
	for _, addr := range listens {
		if addr.Protocol.IsHTTPFamily() {
			httpConfig.IsOn = true
			httpConfig.Listen = append(httpConfig.Listen, addr)
		} else if addr.Protocol.IsHTTPSFamily() {
			httpsConfig.IsOn = true
			httpsConfig.Listen = append(httpsConfig.Listen, addr)
		}
	}

	// 证书
	certIds := []int64{}
	if len(params.CertIdsJSON) > 0 {
		err = json.Unmarshal(params.CertIdsJSON, &certIds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	if httpsConfig.IsOn && len(httpsConfig.Listen) > 0 && len(certIds) == 0 {
		this.Fail("请添加至少一个证书")
	}

	certRefs := []*sslconfigs.SSLCertRef{}
	for _, certId := range certIds {
		certRefs = append(certRefs, &sslconfigs.SSLCertRef{
			IsOn:   true,
			CertId: certId,
		})
	}
	certRefsJSON, err := json.Marshal(certRefs)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建策略
	sslPolicyId := params.SslPolicyId
	if sslPolicyId == 0 {
		if len(certIds) > 0 {
			sslPolicyCreateResp, err := this.RPC().SSLPolicyRPC().CreateSSLPolicy(this.AdminContext(), &pb.CreateSSLPolicyRequest{
				SslCertsJSON: certRefsJSON,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			sslPolicyId = sslPolicyCreateResp.SslPolicyId
		}
	} else {
		_, err = this.RPC().SSLPolicyRPC().UpdateSSLPolicy(this.AdminContext(), &pb.UpdateSSLPolicyRequest{
			SslPolicyId: sslPolicyId,
			SslCertsJSON:   certRefsJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	httpsConfig.SSLPolicyRef = &sslconfigs.SSLPolicyRef{
		IsOn:        true,
		SSLPolicyId: sslPolicyId,
	}

	// 访问地址
	accessAddrs := []*serverconfigs.NetworkAddressConfig{}
	err = json.Unmarshal(params.AccessAddrsJSON, &accessAddrs)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if len(accessAddrs) == 0 {
		this.Fail("请添加至少一个外部访问地址")
	}

	httpJSON, err := json.Marshal(httpConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	httpsJSON, err := json.Marshal(httpsConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().UserNodeRPC().UpdateUserNode(this.AdminContext(), &pb.UpdateUserNodeRequest{
		NodeId:          params.NodeId,
		Name:            params.Name,
		Description:     params.Description,
		HttpJSON:        httpJSON,
		HttpsJSON:       httpsJSON,
		AccessAddrsJSON: params.AccessAddrsJSON,
		IsOn:            params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改用户节点 %d", params.NodeId)

	this.Success()
}
