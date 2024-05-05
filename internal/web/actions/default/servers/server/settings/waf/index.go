package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("waf")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// 只有HTTP服务才支持
	if this.FilterHTTPFamily() {
		return
	}

	// 服务分组设置
	groupResp, err := this.RPC().ServerGroupRPC().FindEnabledServerGroupConfigInfo(this.AdminContext(), &pb.FindEnabledServerGroupConfigInfoRequest{
		ServerId: params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasGroupConfig"] = groupResp.HasWAFConfig
	this.Data["groupSettingURL"] = "/servers/groups/group/settings/waf?groupId=" + types.String(groupResp.ServerGroupId)

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["firewallConfig"] = webConfig.FirewallRef

	// 获取当前网站所在集群的WAF设置
	firewallPolicy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicyWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if firewallPolicy != nil {
		// captcha action
		var captchaOptions = firewallconfigs.NewHTTPFirewallCaptchaAction()
		if len(firewallPolicy.CaptchaOptionsJSON) > 0 {
			err = json.Unmarshal(firewallPolicy.CaptchaOptionsJSON, captchaOptions)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		this.Data["firewallPolicy"] = maps.Map{
			"id":            firewallPolicy.Id,
			"name":          firewallPolicy.Name,
			"isOn":          firewallPolicy.IsOn,
			"mode":          firewallPolicy.Mode,
			"modeInfo":      firewallconfigs.FindFirewallMode(firewallPolicy.Mode),
			"captchaAction": captchaOptions,
		}
	} else {
		this.Data["firewallPolicy"] = nil
	}

	// 当前的Server独立设置
	if webConfig.FirewallRef == nil || webConfig.FirewallRef.FirewallPolicyId == 0 {
		firewallPolicyId, err := dao.SharedHTTPWebDAO.InitEmptyHTTPFirewallPolicy(this.AdminContext(), 0, params.ServerId, webConfig.Id, webConfig.FirewallRef != nil && webConfig.FirewallRef.IsOn)
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
