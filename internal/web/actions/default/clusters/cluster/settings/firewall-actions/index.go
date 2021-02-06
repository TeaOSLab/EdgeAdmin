package firewallActions

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("firewallAction")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	actionsResp, err := this.RPC().NodeClusterFirewallActionRPC().FindAllEnabledNodeClusterFirewallActions(this.AdminContext(), &pb.FindAllEnabledNodeClusterFirewallActionsRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	levelMaps := map[string][]maps.Map{} // level => actionMaps
	for _, action := range actionsResp.NodeClusterFirewallActions {
		actionMaps, ok := levelMaps[action.EventLevel]
		if !ok {
			actionMaps = []maps.Map{}
		}

		actionMaps = append(actionMaps, maps.Map{
			"id":       action.Id,
			"name":     action.Name,
			"type":     action.Type,
			"typeName": firewallconfigs.FindFirewallActionTypeName(action.Type),
		})
		levelMaps[action.EventLevel] = actionMaps
	}

	levelMaps2 := []maps.Map{} // []levelMap
	hasActions := false
	for _, level := range firewallconfigs.FindAllFirewallEventLevels() {
		actionMaps, ok := levelMaps[level.Code]
		if !ok {
			actionMaps = []maps.Map{}
		} else {
			hasActions = true
		}

		levelMaps2 = append(levelMaps2, maps.Map{
			"name":    level.Name,
			"code":    level.Code,
			"actions": actionMaps,
		})
	}

	this.Data["levels"] = levelMaps2
	this.Data["hasActions"] = hasActions

	this.Show()
}
