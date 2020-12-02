package configloaders

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/logs"
	"reflect"
)

var sharedLogConfig *systemconfigs.LogConfig = nil

const (
	LogSettingName = "adminLogConfig"
)

func LoadLogConfig() (*systemconfigs.LogConfig, error) {
	locker.Lock()
	defer locker.Unlock()

	config, err := loadLogConfig()
	if err != nil {
		return nil, err
	}

	v := reflect.Indirect(reflect.ValueOf(config)).Interface().(systemconfigs.LogConfig)
	return &v, nil
}

func UpdateLogConfig(logConfig *systemconfigs.LogConfig) error {
	locker.Lock()
	defer locker.Unlock()

	var rpcClient, err = rpc.SharedRPC()
	if err != nil {
		return err
	}
	valueJSON, err := json.Marshal(logConfig)
	if err != nil {
		return err
	}
	_, err = rpcClient.SysSettingRPC().UpdateSysSetting(rpcClient.Context(0), &pb.UpdateSysSettingRequest{
		Code:      LogSettingName,
		ValueJSON: valueJSON,
	})
	if err != nil {
		return err
	}
	sharedLogConfig = logConfig

	return nil
}

func loadLogConfig() (*systemconfigs.LogConfig, error) {
	if sharedLogConfig != nil {
		return sharedLogConfig, nil
	}
	var rpcClient, err = rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := rpcClient.SysSettingRPC().ReadSysSetting(rpcClient.Context(0), &pb.ReadSysSettingRequest{
		Code: LogSettingName,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.ValueJSON) == 0 {
		sharedLogConfig = systemconfigs.DefaultLogConfig()
		return sharedLogConfig, nil
	}

	config := &systemconfigs.LogConfig{}
	err = json.Unmarshal(resp.ValueJSON, config)
	if err != nil {
		logs.Println("[LOG_MANAGER]" + err.Error())
		sharedLogConfig = systemconfigs.DefaultLogConfig()
		return sharedLogConfig, nil
	}
	sharedLogConfig = config
	return sharedLogConfig, nil
}
