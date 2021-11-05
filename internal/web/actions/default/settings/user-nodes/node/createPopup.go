package node

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
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

	createResp, err := this.RPC().UserNodeRPC().CreateUserNode(this.AdminContext(), &pb.CreateUserNodeRequest{
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
	defer this.CreateLog(oplogs.LevelInfo, "创建用户节点 %d", createResp.UserNodeId)

	this.Success()
}
