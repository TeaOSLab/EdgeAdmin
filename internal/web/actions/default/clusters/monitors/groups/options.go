// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type OptionsAction struct {
	actionutils.ParentAction
}

func (this *OptionsAction) RunPost(params struct{}) {
	resp, err := this.RPC().ReportNodeGroupRPC().FindAllEnabledReportNodeGroups(this.AdminContext(), &pb.FindAllEnabledReportNodeGroupsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var groupMaps = []maps.Map{}
	for _, group := range resp.ReportNodeGroups {
		groupMaps = append(groupMaps, maps.Map{
			"id":   group.Id,
			"name": group.Name,
		})
	}
	this.Data["groups"] = groupMaps

	this.Success()
}
