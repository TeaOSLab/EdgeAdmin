// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package logs

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/actions"
)

type SettingsAction struct {
	actionutils.ParentAction
}

func (this *SettingsAction) Init() {
	this.Nav("", "", "settings")
}

func (this *SettingsAction) RunGet(params struct{}) {
	settingsResp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeAccessLogQueue})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var config = &serverconfigs.AccessLogQueueConfig{
		MaxLength:         0,
		CountPerSecond:    0,
		Percent:           100,
		EnableAutoPartial: true,
		RowsPerTable:      500_000,
	}
	if len(settingsResp.ValueJSON) > 0 {
		err = json.Unmarshal(settingsResp.ValueJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		configJSON, err := json.Marshal(config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().SysSettingRPC().UpdateSysSetting(this.AdminContext(), &pb.UpdateSysSettingRequest{
			Code:      systemconfigs.SettingCodeAccessLogQueue,
			ValueJSON: configJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["config"] = config

	this.Show()
}

func (this *SettingsAction) RunPost(params struct {
	Percent           int
	CountPerSecond    int
	MaxLength         int
	EnableAutoPartial bool
	RowsPerTable      int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("percent", params.Percent).
		Gte(0, "请输入大于0的整数").
		Lte(100, "请输入小于100的整数")

	var config = &serverconfigs.AccessLogQueueConfig{
		MaxLength:         params.MaxLength,
		CountPerSecond:    params.CountPerSecond,
		Percent:           params.Percent,
		EnableAutoPartial: params.EnableAutoPartial,
		RowsPerTable:      params.RowsPerTable,
	}
	configJSON, err := json.Marshal(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().SysSettingRPC().UpdateSysSetting(this.AdminContext(), &pb.UpdateSysSettingRequest{
		Code:      systemconfigs.SettingCodeAccessLogQueue,
		ValueJSON: configJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
