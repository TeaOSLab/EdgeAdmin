package node

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "node", "create")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name            string
	Description     string
	ListensJSON     []byte
	CertIdsJSON     []byte
	AccessAddrsJSON []byte

	RestIsOn        bool
	RestListensJSON []byte

	IsOn bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入API节点名称")

	var httpConfig = &serverconfigs.HTTPProtocolConfig{}
	var httpsConfig = &serverconfigs.HTTPSProtocolConfig{}

	// 监听地址
	var listens = []*serverconfigs.NetworkAddressConfig{}
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

	// Rest监听地址
	var restHTTPConfig = &serverconfigs.HTTPProtocolConfig{}
	var restHTTPSConfig = &serverconfigs.HTTPSProtocolConfig{}
	if params.RestIsOn {
		var restListens = []*serverconfigs.NetworkAddressConfig{}
		err = json.Unmarshal(params.RestListensJSON, &restListens)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(restListens) == 0 {
			this.Fail("请至少添加一个HTTP API监听端口")
			return
		}
		for _, addr := range restListens {
			if addr.Protocol.IsHTTPFamily() {
				restHTTPConfig.IsOn = true
				restHTTPConfig.Listen = append(restHTTPConfig.Listen, addr)
			} else if addr.Protocol.IsHTTPSFamily() {
				restHTTPSConfig.IsOn = true
				restHTTPSConfig.Listen = append(restHTTPSConfig.Listen, addr)
			}
		}

		// 是否有端口冲突
		var rpcAddresses = []string{}
		for _, listen := range listens {
			err := listen.Init()
			if err != nil {
				this.Fail("校验配置失败：" + configutils.QuoteIP(listen.Host) + ":" + listen.PortRange + ": " + err.Error())
				return
			}
			rpcAddresses = append(rpcAddresses, listen.Addresses()...)
		}

		for _, listen := range restListens {
			err := listen.Init()
			if err != nil {
				this.Fail("校验配置失败：" + configutils.QuoteIP(listen.Host) + ":" + listen.PortRange + ": " + err.Error())
				return
			}
			for _, address := range listen.Addresses() {
				if lists.ContainsString(rpcAddresses, address) {
					this.Fail("HTTP API地址 '" + address + "' 和 GRPC地址冲突，请修改后提交")
					return
				}
			}
		}
	}

	// 证书
	var certIds = []int64{}
	if len(params.CertIdsJSON) > 0 {
		err = json.Unmarshal(params.CertIdsJSON, &certIds)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	if ((httpsConfig.IsOn && len(httpsConfig.Listen) > 0) || (restHTTPSConfig.IsOn && len(httpsConfig.Listen) > 0)) && len(certIds) == 0 {
		this.Fail("请添加至少一个证书")
	}

	var certRefs = []*sslconfigs.SSLCertRef{}
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
	if len(certIds) > 0 {
		sslPolicyCreateResp, err := this.RPC().SSLPolicyRPC().CreateSSLPolicy(this.AdminContext(), &pb.CreateSSLPolicyRequest{
			SslCertsJSON: certRefsJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		sslPolicyId := sslPolicyCreateResp.SslPolicyId
		httpsConfig.SSLPolicyRef = &sslconfigs.SSLPolicyRef{
			IsOn:        true,
			SSLPolicyId: sslPolicyId,
		}
		restHTTPSConfig.SSLPolicyRef = &sslconfigs.SSLPolicyRef{
			IsOn:        true,
			SSLPolicyId: sslPolicyId,
		}
	}

	// 访问地址
	var accessAddrs = []*serverconfigs.NetworkAddressConfig{}
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

	restHTTPJSON, err := json.Marshal(restHTTPConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	restHTTPSJSON, err := json.Marshal(restHTTPSConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	createResp, err := this.RPC().APINodeRPC().CreateAPINode(this.AdminContext(), &pb.CreateAPINodeRequest{
		Name:            params.Name,
		Description:     params.Description,
		HttpJSON:        httpJSON,
		HttpsJSON:       httpsJSON,
		RestIsOn:        params.RestIsOn,
		RestHTTPJSON:    restHTTPJSON,
		RestHTTPSJSON:   restHTTPSJSON,
		AccessAddrsJSON: params.AccessAddrsJSON,
		IsOn:            params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLogInfo(codes.APINode_LogCreateAPINode, createResp.ApiNodeId)

	this.Success()
}
