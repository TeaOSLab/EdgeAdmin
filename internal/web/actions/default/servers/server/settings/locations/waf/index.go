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
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
}

func (this *IndexAction) RunGet(params struct {
	ServerId   int64
	LocationId int64
}) {
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithLocationId(this.AdminContext(), params.LocationId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["firewallConfig"] = webConfig.FirewallRef

	// 获取当前服务所在集群的WAF设置
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
			"id":             firewallPolicy.Id,
			"name":           firewallPolicy.Name,
			"isOn":           firewallPolicy.IsOn,
			"mode":           firewallPolicy.Mode,
			"modeInfo":       firewallconfigs.FindFirewallMode(firewallPolicy.Mode),
			"captchaOptions": captchaOptions,
		}
	} else {
		this.Data["firewallPolicy"] = nil
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
