// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package users

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type PlansAction struct {
	actionutils.ParentAction
}

func (this *PlansAction) RunPost(params struct {
	UserId   int64
	ServerId int64
}) {
	if !teaconst.IsPlus || params.UserId <= 0 {
		this.Data["plans"] = []maps.Map{}
		this.Success()
	}

	// TODO 优化用户套餐查询
	userPlansResp, err := this.RPC().UserPlanRPC().FindAllEnabledUserPlansForServer(this.AdminContext(), &pb.FindAllEnabledUserPlansForServerRequest{
		UserId:   params.UserId,
		ServerId: params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var userPlanMaps = []maps.Map{}
	for _, userPlan := range userPlansResp.UserPlans {
		if userPlan.Plan == nil {
			continue
		}

		userPlanMaps = append(userPlanMaps, maps.Map{
			"id":    userPlan.Id,
			"name":  userPlan.Plan.Name,
			"dayTo": userPlan.DayTo,
		})
	}
	this.Data["plans"] = userPlanMaps

	this.Success()
}
