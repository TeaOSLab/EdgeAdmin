package waf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/groups/group/servergrouputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("waf")
}

func (this *IndexAction) RunGet(params struct {
	GroupId int64
}) {
	_, err := servergrouputils.InitGroup(this.Parent(), params.GroupId, "waf")
	if err != nil {
		this.ErrorPage(err)
		return
	}

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerGroupId(this.AdminContext(), params.GroupId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["firewallConfig"] = webConfig.FirewallRef

	// 获取当前服务所在集群的WAF设置
	this.Data["firewallPolicy"] = nil

	// 当前的Server独立设置
	if webConfig.FirewallRef == nil || webConfig.FirewallRef.FirewallPolicyId == 0 {
		firewallPolicyId, err := dao.SharedHTTPWebDAO.InitEmptyHTTPFirewallPolicy(this.AdminContext(), params.GroupId, 0, webConfig.Id, webConfig.FirewallRef != nil && webConfig.FirewallRef.IsOn)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.Data["firewallPolicyId"] = firewallPolicyId
	} else {
		this.Data["firewallPolicyId"] = webConfig.FirewallRef.FirewallPolicyId
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId        int64
	FirewallJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.ServerWAF_LogUpdateWAFSettings, params.WebId)

	// TODO 检查配置

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebFirewall(this.AdminContext(), &pb.UpdateHTTPWebFirewallRequest{
		HttpWebId:    params.WebId,
		FirewallJSON: params.FirewallJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
