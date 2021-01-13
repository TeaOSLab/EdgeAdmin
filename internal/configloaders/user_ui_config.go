package configloaders

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/logs"
	"reflect"
)

var sharedUserUIConfig *systemconfigs.UserUIConfig = nil

const (
	UserUISettingName = "userUIConfig"
)

func LoadUserUIConfig() (*systemconfigs.UserUIConfig, error) {
	locker.Lock()
	defer locker.Unlock()

	config, err := loadUserUIConfig()
	if err != nil {
		return nil, err
	}

	v := reflect.Indirect(reflect.ValueOf(config)).Interface().(systemconfigs.UserUIConfig)
	return &v, nil
}

func UpdateUserUIConfig(uiConfig *systemconfigs.UserUIConfig) error {
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
		Code:      UserUISettingName,
		ValueJSON: valueJSON,
	})
	if err != nil {
		return err
	}
	sharedUserUIConfig = uiConfig

	return nil
}

func loadUserUIConfig() (*systemconfigs.UserUIConfig, error) {
	if sharedUserUIConfig != nil {
		return sharedUserUIConfig, nil
	}
	var rpcClient, err = rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := rpcClient.SysSettingRPC().ReadSysSetting(rpcClient.Context(0), &pb.ReadSysSettingRequest{
		Code: UserUISettingName,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.ValueJSON) == 0 {
		sharedUserUIConfig = defaultUserUIConfig()
		return sharedUserUIConfig, nil
	}

	config := &systemconfigs.UserUIConfig{}
	err = json.Unmarshal(resp.ValueJSON, config)
	if err != nil {
		logs.Println("[UI_MANAGER]" + err.Error())
		sharedUserUIConfig = defaultUserUIConfig()
		return sharedUserUIConfig, nil
	}
	sharedUserUIConfig = config
	return sharedUserUIConfig, nil
}

func defaultUserUIConfig() *systemconfigs.UserUIConfig {
	return &systemconfigs.UserUIConfig{
		ProductName:        "GoEdge",
		UserSystemName:     "GoEdge用户系统",
		ShowOpenSourceInfo: true,
		ShowVersion:        true,
		ShowFinance:        true,
	}
}
