package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
)

type SortSetsAction struct {
	actionutils.ParentAction
}

func (this *SortSetsAction) RunPost(params struct {
	GroupId int64
	SetIds  []int64
}) {
	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "修改WAF规则分组 %d 中的规则集排序", params.GroupId)

	groupConfig, err := dao.SharedHTTPFirewallRuleGroupDAO.FindRuleGroupConfig(this.AdminContext(), params.GroupId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if groupConfig == nil {
		this.NotFound("firewallRuleGroup", params.GroupId)
		return
	}

	setMap := map[int64]*firewallconfigs.HTTPFirewallRuleSetRef{}
	for _, setRef := range groupConfig.SetRefs {
		setMap[setRef.SetId] = setRef
	}

	newRefs := []*firewallconfigs.HTTPFirewallRuleSetRef{}
	for _, setId := range params.SetIds {
		ref, ok := setMap[setId]
		if ok {
			newRefs = append(newRefs, ref)
		}
	}
	newRefsJSON, err := json.Marshal(newRefs)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPFirewallRuleGroupRPC().UpdateHTTPFirewallRuleGroupSets(this.AdminContext(), &pb.UpdateHTTPFirewallRuleGroupSetsRequest{
		FirewallRuleGroupId:  params.GroupId,
		FirewallRuleSetsJSON: newRefsJSON,
	})

	this.Success()
}
