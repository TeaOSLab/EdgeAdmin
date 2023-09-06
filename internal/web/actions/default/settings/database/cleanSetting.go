package database

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
	var config = systemconfigs.NewDatabaseConfig()
	if len(configResp.ValueJSON) > 0 {
		err = json.Unmarshal(configResp.ValueJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["config"] = config

	this.Show()
}

func (this *CleanSettingAction) RunPost(params struct {
	ServerAccessLogCleanDays             int
	ServerBandwidthStatCleanDays         int
	UserBandwidthStatCleanDays           int
	UserPlanBandwidthStatCleanDays       int
	ServerDailyStatCleanDays             int
	ServerDomainHourlyStatCleanDays      int
	TrafficDailyStatCleanDays            int
	TrafficHourlyStatCleanDays           int
	NodeClusterTrafficDailyStatCleanDays int
	NodeTrafficDailyStatCleanDays        int
	NodeTrafficHourlyStatCleanDays       int
	HttpCacheTaskCleanDays               int

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.Database_LogUpdateCleanDays)

	// 读取设置
	configResp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeDatabaseConfigSetting})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var config = systemconfigs.NewDatabaseConfig()
	if len(configResp.ValueJSON) > 0 {
		err = json.Unmarshal(configResp.ValueJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	if params.ServerAccessLogCleanDays < 0 {
		params.ServerAccessLogCleanDays = 0
	}
	config.ServerAccessLog.Clean.Days = params.ServerAccessLogCleanDays

	if params.ServerBandwidthStatCleanDays < 0 {
		params.ServerBandwidthStatCleanDays = 0
	}
	config.ServerBandwidthStat.Clean.Days = params.ServerBandwidthStatCleanDays

	if params.UserBandwidthStatCleanDays < 0 {
		params.UserBandwidthStatCleanDays = 0
	}
	config.UserBandwidthStat.Clean.Days = params.UserBandwidthStatCleanDays

	if params.UserPlanBandwidthStatCleanDays < 0 {
		params.UserPlanBandwidthStatCleanDays = 0
	}
	config.UserPlanBandwidthStat.Clean.Days = params.UserPlanBandwidthStatCleanDays

	if params.ServerDailyStatCleanDays < 0 {
		params.ServerDailyStatCleanDays = 0
	}
	config.ServerDailyStat.Clean.Days = params.ServerDailyStatCleanDays

	if params.ServerDomainHourlyStatCleanDays < 0 {
		params.ServerDomainHourlyStatCleanDays = 0
	}
	config.ServerDomainHourlyStat.Clean.Days = params.ServerDomainHourlyStatCleanDays

	if params.TrafficDailyStatCleanDays < 0 {
		params.TrafficDailyStatCleanDays = 0
	}
	config.TrafficDailyStat.Clean.Days = params.TrafficDailyStatCleanDays

	if params.TrafficHourlyStatCleanDays < 0 {
		params.TrafficHourlyStatCleanDays = 0
	}
	config.TrafficHourlyStat.Clean.Days = params.TrafficHourlyStatCleanDays

	if params.NodeClusterTrafficDailyStatCleanDays < 0 {
		params.NodeClusterTrafficDailyStatCleanDays = 0
	}
	config.NodeClusterTrafficDailyStat.Clean.Days = params.NodeClusterTrafficDailyStatCleanDays

	if params.NodeTrafficDailyStatCleanDays < 0 {
		params.NodeTrafficDailyStatCleanDays = 0
	}
	config.NodeTrafficDailyStat.Clean.Days = params.NodeTrafficDailyStatCleanDays

	if params.NodeTrafficHourlyStatCleanDays < 0 {
		params.NodeTrafficHourlyStatCleanDays = 0
	}
	config.NodeTrafficHourlyStat.Clean.Days = params.NodeTrafficHourlyStatCleanDays

	if params.HttpCacheTaskCleanDays < 0 {
		params.HttpCacheTaskCleanDays = 0
	}
	config.HTTPCacheTask.Clean.Days = params.HttpCacheTaskCleanDays

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
