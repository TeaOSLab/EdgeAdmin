package securitymanager

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/logs"
	"reflect"
	"sync"
)

var locker sync.Mutex

const (
	SecuritySettingName = "adminSecurityConfig"

	FrameNone       = ""
	FrameDeny       = "DENY"
	FrameSameOrigin = "SAMEORIGIN"
)

var sharedSecurityConfig *SecurityConfig = nil

type SecurityConfig struct {
	Frame            string  `json:"frame"`
	AllowCountryIds  []int64 `json:"allowCountryIds"`
	AllowProvinceIds []int64 `json:"allowProvinceIds"`
	AllowLocal       bool    `json:"allowLocal"`
}

func LoadSecurityConfig() (*SecurityConfig, error) {
	locker.Lock()
	defer locker.Unlock()

	config, err := loadSecurityConfig()
	if err != nil {
		return nil, err
	}

	v := reflect.Indirect(reflect.ValueOf(config)).Interface().(SecurityConfig)
	return &v, nil
}

func UpdateSecurityConfig(securityConfig *SecurityConfig) error {
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
	sharedSecurityConfig = securityConfig

	// 通知更新
	events.Notify(events.EventSecurityConfigChanged)

	return nil
}

func loadSecurityConfig() (*SecurityConfig, error) {
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
		sharedSecurityConfig = defaultSecurityConfig()
		return sharedSecurityConfig, nil
	}

	config := &SecurityConfig{}
	err = json.Unmarshal(resp.ValueJSON, config)
	if err != nil {
		logs.Println("[SECURITY_MANAGER]" + err.Error())
		sharedSecurityConfig = defaultSecurityConfig()
		return sharedSecurityConfig, nil
	}
	sharedSecurityConfig = config
	return sharedSecurityConfig, nil
}

func defaultSecurityConfig() *SecurityConfig {
	return &SecurityConfig{
		Frame:      FrameSameOrigin,
		AllowLocal: true,
	}
}
