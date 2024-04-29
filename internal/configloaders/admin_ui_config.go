package configloaders

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/logs"
	"reflect"
	"time"
)

var sharedAdminUIConfig *systemconfigs.AdminUIConfig = nil

func LoadAdminUIConfig() (*systemconfigs.AdminUIConfig, error) {
	locker.Lock()
	defer locker.Unlock()

	config, err := loadAdminUIConfig()
	if err != nil {
		return nil, err
	}

	v := reflect.Indirect(reflect.ValueOf(config)).Interface().(systemconfigs.AdminUIConfig)
	return &v, nil
}

func ReloadAdminUIConfig() error {
	locker.Lock()
	defer locker.Unlock()

	sharedAdminUIConfig = nil
	_, err := loadAdminUIConfig()
	return err
}

func UpdateAdminUIConfig(uiConfig *systemconfigs.AdminUIConfig) error {
	locker.Lock()
	defer locker.Unlock()

	var rpcClient, err = rpc.SharedRPC()
	if err != nil {
		return err
	}
	valueJSON, err := json.Marshal(uiConfig)
	if err != nil {
		return err
	}
	_, err = rpcClient.SysSettingRPC().UpdateSysSetting(rpcClient.Context(0), &pb.UpdateSysSettingRequest{
		Code:      systemconfigs.SettingCodeAdminUIConfig,
		ValueJSON: valueJSON,
	})
	if err != nil {
		return err
	}
	sharedAdminUIConfig = uiConfig

	// timezone
	updateTimeZone(uiConfig)

	return nil
}

// ShowFinance 是否显示财务信息
func ShowFinance() bool {
	config, _ := LoadAdminUIConfig()
	if config != nil && !config.ShowFinance {
		return false
	}
	return true
}

func loadAdminUIConfig() (*systemconfigs.AdminUIConfig, error) {
	if sharedAdminUIConfig != nil {
		return sharedAdminUIConfig, nil
	}
	var rpcClient, err = rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := rpcClient.SysSettingRPC().ReadSysSetting(rpcClient.Context(0), &pb.ReadSysSettingRequest{
		Code: systemconfigs.SettingCodeAdminUIConfig,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.ValueJSON) == 0 {
		sharedAdminUIConfig = defaultAdminUIConfig()
		return sharedAdminUIConfig, nil
	}

	var config = &systemconfigs.AdminUIConfig{}
	config.DNSResolver.Type = nodeconfigs.DNSResolverTypeDefault // 默认值
	err = json.Unmarshal(resp.ValueJSON, config)
	if err != nil {
		logs.Println("[UI_MANAGER]" + err.Error())
		sharedAdminUIConfig = defaultAdminUIConfig()
		return sharedAdminUIConfig, nil
	}

	// timezone
	updateTimeZone(config)

	sharedAdminUIConfig = config
	return sharedAdminUIConfig, nil
}

func defaultAdminUIConfig() *systemconfigs.AdminUIConfig {
	var config = &systemconfigs.AdminUIConfig{
		ProductName:        langs.DefaultMessage(codes.AdminUI_DefaultProductName),
		AdminSystemName:    langs.DefaultMessage(codes.AdminUI_DefaultSystemName),
		ShowOpenSourceInfo: true,
		ShowVersion:        true,
		ShowFinance:        true,
		DefaultPageSize:    10,
		TimeZone:           nodeconfigs.DefaultTimeZoneLocation,
	}
	config.DNSResolver.Type = nodeconfigs.DNSResolverTypeDefault
	return config
}

// 修改时区
func updateTimeZone(config *systemconfigs.AdminUIConfig) {
	if len(config.TimeZone) > 0 {
		location, err := time.LoadLocation(config.TimeZone)
		if err == nil && time.Local != location {
			time.Local = location
		}
	}
}
