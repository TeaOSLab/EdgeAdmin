package uimanager

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/logs"
	"reflect"
	"sync"
)

var sharedUIConfig *UIConfig = nil
var locker sync.Mutex

const (
	UISettingName = "adminUIConfig"
)

type UIConfig struct {
	ProductName        string `json:"productName"`        // 产品名
	AdminSystemName    string `json:"adminSystemName"`    // 管理员系统名称
	ShowOpenSourceInfo bool   `json:"showOpenSourceInfo"` // 是否显示开源信息
	ShowVersion        bool   `json:"showVersion"`        // 是否显示版本号
	Version            string `json:"version"`            // 显示的版本号
}

func LoadUIConfig() (*UIConfig, error) {
	locker.Lock()
	defer locker.Unlock()

	config, err := loadUIConfig()
	if err != nil {
		return nil, err
	}

	v := reflect.Indirect(reflect.ValueOf(config)).Interface().(UIConfig)
	return &v, nil
}

func UpdateUIConfig(uiConfig *UIConfig) error {
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

func loadUIConfig() (*UIConfig, error) {
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

	config := &UIConfig{}
	err = json.Unmarshal(resp.ValueJSON, config)
	if err != nil {
		logs.Println("[UI_MANAGER]" + err.Error())
		sharedUIConfig = defaultUIConfig()
		return sharedUIConfig, nil
	}
	sharedUIConfig = config
	return sharedUIConfig, nil
}

func defaultUIConfig() *UIConfig {
	return &UIConfig{
		ProductName:        "GoEdge",
		AdminSystemName:    "GoEdge管理员系统",
		ShowOpenSourceInfo: true,
		ShowVersion:        true,
	}
}
