package configloaders

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/logs"
	"reflect"
)

var sharedUIConfig *systemconfigs.AdminUIConfig = nil

const (
	UISettingName = "adminUIConfig"
)

func LoadUIConfig() (*systemconfigs.AdminUIConfig, error) {
	locker.Lock()
	defer locker.Unlock()

	config, err := loadUIConfig()
	if err != nil {
		return nil, err
	}

	v := reflect.Indirect(reflect.ValueOf(config)).Interface().(systemconfigs.AdminUIConfig)
	return &v, nil
}

func UpdateUIConfig(uiConfig *systemconfigs.AdminUIConfig) error {
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
		Code:      UISettingName,
		ValueJSON: valueJSON,
	})
	if err != nil {
		return err
	}
	sharedUIConfig = uiConfig

	return nil
}

func loadUIConfig() (*systemconfigs.AdminUIConfig, error) {
	if sharedUIConfig != nil {
		return sharedUIConfig, nil
	}
	var rpcClient, err = rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := rpcClient.SysSettingRPC().ReadSysSetting(rpcClient.Context(0), &pb.ReadSysSettingRequest{
		Code: UISettingName,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.ValueJSON) == 0 {
		sharedUIConfig = defaultUIConfig()
		return sharedUIConfig, nil
	}

	config := &systemconfigs.AdminUIConfig{}
	err = json.Unmarshal(resp.ValueJSON, config)
	if err != nil {
		logs.Println("[UI_MANAGER]" + err.Error())
		sharedUIConfig = defaultUIConfig()
		return sharedUIConfig, nil
	}
	sharedUIConfig = config
	return sharedUIConfig, nil
}

func defaultUIConfig() *systemconfigs.AdminUIConfig {
	return &systemconfigs.AdminUIConfig{
		ProductName:        "GoEdge",
		AdminSystemName:    "GoEdge管理员系统",
		ShowOpenSourceInfo: true,
		ShowVersion:        true,
	}
}
