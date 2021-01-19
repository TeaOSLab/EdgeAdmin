package database

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/actions"
)

type CleanSettingAction struct {
	actionutils.ParentAction
}

func (this *CleanSettingAction) Init() {
	this.Nav("", "", "cleanSetting")
}

func (this *CleanSettingAction) RunGet(params struct{}) {
	// 读取设置
	configResp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeDatabaseConfigSetting})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var config = &systemconfigs.DatabaseConfig{}
	if len(configResp.ValueJSON) > 0 {
		err = json.Unmarshal(configResp.ValueJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["config"] = config.ServerAccessLog.Clean

	this.Show()
}

func (this *CleanSettingAction) RunPost(params struct {
	Days int

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改数据库自动清理设置")

	days := params.Days
	if days < 0 {
		days = 0
	}

	// 读取设置
	configResp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeDatabaseConfigSetting})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var config = &systemconfigs.DatabaseConfig{}
	if len(configResp.ValueJSON) > 0 {
		err = json.Unmarshal(configResp.ValueJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	config.ServerAccessLog.Clean.Days = days
	configJSON, err := json.Marshal(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().SysSettingRPC().UpdateSysSetting(this.AdminContext(), &pb.UpdateSysSettingRequest{
		Code:      systemconfigs.SettingCodeDatabaseConfigSetting,
		ValueJSON: configJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
