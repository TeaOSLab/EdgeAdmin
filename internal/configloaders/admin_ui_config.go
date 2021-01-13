package configloaders

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/logs"
	"reflect"
)

var sharedAdminUIConfig *systemconfigs.AdminUIConfig = nil

const (
	AdminUISettingName = "adminUIConfig"
)

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
		Code:      AdminUISettingName,
		ValueJSON: valueJSON,
	})
	if err != nil {
		return err
	}
	sharedAdminUIConfig = uiConfig

	return nil
}

// 是否显示财务信息
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
		Code: AdminUISettingName,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.ValueJSON) == 0 {
		sharedAdminUIConfig = defaultAdminUIConfig()
		return sharedAdminUIConfig, nil
	}

	config := &systemconfigs.AdminUIConfig{}
	err = json.Unmarshal(resp.ValueJSON, config)
	if err != nil {
		logs.Println("[UI_MANAGER]" + err.Error())
		sharedAdminUIConfig = defaultAdminUIConfig()
		return sharedAdminUIConfig, nil
	}
	sharedAdminUIConfig = config
	return sharedAdminUIConfig, nil
}

func defaultAdminUIConfig() *systemconfigs.AdminUIConfig {
	return &systemconfigs.AdminUIConfig{
		ProductName:        "GoEdge",
		AdminSystemName:    "GoEdge管理员系统",
		ShowOpenSourceInfo: true,
		ShowVersion:        true,
		ShowFinance:        true,
	}
}
