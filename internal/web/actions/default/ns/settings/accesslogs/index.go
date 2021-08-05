// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package accesslogs

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
	this.SecondMenu("accessLog")
}

func (this *IndexAction) RunGet(params struct{}) {
	resp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{
		Code: systemconfigs.SettingCodeNSAccessLogSetting,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var config = &dnsconfigs.NSAccessLogRef{}
	if len(resp.ValueJSON) > 0 {
		err = json.Unmarshal(resp.ValueJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		config.IsOn = true
		config.LogMissingDomains = true
	}

	this.Data["config"] = config

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	AccessLogJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// 校验配置
	var config = &dnsconfigs.NSAccessLogRef{}
	err := json.Unmarshal(params.AccessLogJSON, config)
	if err != nil {
		this.Fail("配置解析失败：" + err.Error())
	}

	err = config.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
	}

	_, err = this.RPC().SysSettingRPC().UpdateSysSetting(this.AdminContext(), &pb.UpdateSysSettingRequest{
		Code:      systemconfigs.SettingCodeNSAccessLogSetting,
		ValueJSON: params.AccessLogJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
