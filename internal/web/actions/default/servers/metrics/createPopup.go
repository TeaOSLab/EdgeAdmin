// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package metrics

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
	Category string
}) {
	this.Data["category"] = params.Category
	this.Data["valueDefinitions"] = serverconfigs.FindAllMetricValueDefinitions(params.Category)

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name          string
	Category      string
	KeysJSON      []byte
	PeriodJSON    []byte
	Value         string
	IsPublic      bool
	ExpiresPeriod int32

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入指标名称")

	if len(params.Category) == 0 {
		this.Fail("请选择指标类型")
	}

	// 统计对象
	if len(params.KeysJSON) == 0 {
		this.FailField("keys", "请选择指标统计的对象")
	}
	var keys = []string{}
	err := json.Unmarshal(params.KeysJSON, &keys)
	if err != nil {
		this.FailField("keys", "解析指标对象失败")
	}
	if len(keys) == 0 {
		this.FailField("keys", "请选择指标统计的对象")
	}

	var periodMap = maps.Map{}
	err = json.Unmarshal(params.PeriodJSON, &periodMap)
	if err != nil {
		this.FailField("period", "解析统计周期失败")
	}
	var period = periodMap.GetInt32("period")
	var periodUnit = periodMap.GetString("unit")

	if params.ExpiresPeriod < 0 {
		params.ExpiresPeriod = 0
	}

	createResp, err := this.RPC().MetricItemRPC().CreateMetricItem(this.AdminContext(), &pb.CreateMetricItemRequest{
		Code:          "", // TODO 未来实现
		Category:      params.Category,
		Name:          params.Name,
		Keys:          keys,
		Period:        period,
		PeriodUnit:    periodUnit,
		ExpiresPeriod: params.ExpiresPeriod,
		Value:         params.Value,
		IsPublic:      params.IsPublic,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	defer this.CreateLogInfo(codes.MetricItem_LogCreateMetricItem, createResp.MetricItemId)
	this.Success()
}
