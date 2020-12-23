package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type Helper struct {
}

func NewHelper() *Helper {
	return &Helper{}
}

func (this *Helper) BeforeAction(actionPtr actions.ActionWrapper) (goNext bool) {
	action := actionPtr.Object()
	if action.Request.Method != http.MethodGet {
		return true
	}

	action.Data["mainTab"] = "component"
	action.Data["secondMenuItem"] = "waf"

	// 显示当前的防火墙名称
	firewallPolicyId := action.ParamInt64("firewallPolicyId")
	if firewallPolicyId > 0 {
		action.Data["firewallPolicyId"] = firewallPolicyId
		action.Data["countInboundGroups"] = 0
		action.Data["countOutboundGroups"] = 0
		parentAction := actionutils.FindParentAction(actionPtr)
		if parentAction != nil {
			firewallPolicy, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledHTTPFirewallPolicy(parentAction.AdminContext(), firewallPolicyId)
			if err != nil {
				parentAction.ErrorPage(err)
				return
			}
			if firewallPolicy == nil {
				action.WriteString("can not find firewall policy")
				return
			}
			action.Data["firewallPolicyName"] = firewallPolicy.Name

			// inbound
			if len(firewallPolicy.InboundJSON) > 0 {
				inboundConfig := &firewallconfigs.HTTPFirewallInboundConfig{}
				err = json.Unmarshal(firewallPolicy.InboundJSON, inboundConfig)
				if err != nil {
					parentAction.ErrorPage(err)
					return
				}
				action.Data["countInboundGroups"] = len(inboundConfig.GroupRefs)
			}

			// outbound
			if len(firewallPolicy.OutboundJSON) > 0 {
				outboundConfig := &firewallconfigs.HTTPFirewallOutboundConfig{}
				err = json.Unmarshal(firewallPolicy.OutboundJSON, outboundConfig)
				if err != nil {
					parentAction.ErrorPage(err)
					return
				}
				action.Data["countOutboundGroups"] = len(outboundConfig.GroupRefs)
			}
		}
	}
	return true
}
