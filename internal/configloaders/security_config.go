package configloaders

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/logs"
	"reflect"
)

const (
	SecuritySettingName = "adminSecurityConfig"

	FrameNone       = ""
	FrameDeny       = "DENY"
	FrameSameOrigin = "SAMEORIGIN"
)

var sharedSecurityConfig *systemconfigs.SecurityConfig = nil

func LoadSecurityConfig() (*systemconfigs.SecurityConfig, error) {
	locker.Lock()
	defer locker.Unlock()

	config, err := loadSecurityConfig()
	if err != nil {
		return nil, err
	}

	var v = reflect.Indirect(reflect.ValueOf(config)).Interface().(systemconfigs.SecurityConfig)
	return &v, nil
}

func UpdateSecurityConfig(securityConfig *systemconfigs.SecurityConfig) error {
	locker.Lock()
	defer locker.Unlock()

	var rpcClient, err = rpc.SharedRPC()
	if err != nil {
		return err
	}
	valueJSON, err := json.Marshal(securityConfig)
	if err != nil {
		return err
	}
	_, err = rpcClient.SysSettingRPC().UpdateSysSetting(rpcClient.Context(0), &pb.UpdateSysSettingRequest{
		Code:      SecuritySettingName,
		ValueJSON: valueJSON,
	})
	if err != nil {
		return err
	}
	err = securityConfig.Init()
	if err != nil {
		return err
	}
	sharedSecurityConfig = securityConfig

	// 通知更新
	events.Notify(events.EventSecurityConfigChanged)

	return nil
}

func loadSecurityConfig() (*systemconfigs.SecurityConfig, error) {
	if sharedSecurityConfig != nil {
		return sharedSecurityConfig, nil
	}
	var rpcClient, err = rpc.SharedRPC()
	if err != nil {
		return nil, err
	}
	resp, err := rpcClient.SysSettingRPC().ReadSysSetting(rpcClient.Context(0), &pb.ReadSysSettingRequest{
		Code: SecuritySettingName,
	})
	if err != nil {
		return nil, err
	}
	if len(resp.ValueJSON) == 0 {
		sharedSecurityConfig = NewSecurityConfig()
		return sharedSecurityConfig, nil
	}

	var config = &systemconfigs.SecurityConfig{
		Frame:                  FrameSameOrigin,
		AllowLocal:             true,
		CheckClientFingerprint: false,
		CheckClientRegion:      true,
		DenySearchEngines:      true,
		DenySpiders:            true,
	}
	err = json.Unmarshal(resp.ValueJSON, config)
	if err != nil {
		logs.Println("[SECURITY_MANAGER]" + err.Error())
		sharedSecurityConfig = NewSecurityConfig()
		return sharedSecurityConfig, nil
	}
	err = config.Init()
	if err != nil {
		return nil, err
	}
	sharedSecurityConfig = config
	return sharedSecurityConfig, nil
}

// NewSecurityConfig create new security config
func NewSecurityConfig() *systemconfigs.SecurityConfig {
	return &systemconfigs.SecurityConfig{
		Frame:                  FrameSameOrigin,
		AllowLocal:             true,
		CheckClientFingerprint: false,
		CheckClientRegion:      true,
		DenySearchEngines:      true,
		DenySpiders:            true,
	}
}
