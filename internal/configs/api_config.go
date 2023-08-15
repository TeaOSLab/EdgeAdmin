package configs

import (
	"errors"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/Tea"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const ConfigFileName = "api_admin.yaml"
const oldConfigFileName = "api.yaml"

// APIConfig API配置
type APIConfig struct {
	OldRPC struct {
		Endpoints     []string `yaml:"endpoints"`
		DisableUpdate bool     `yaml:"disableUpdate"`
	} `yaml:"rpc,omitempty"`

	RPCEndpoints     []string `yaml:"rpc.endpoints,flow" json:"rpc.endpoints"`
	RPCDisableUpdate bool     `yaml:"rpc.disableUpdate" json:"rpc.disableUpdate"`

	NodeId string `yaml:"nodeId"`
	Secret string `yaml:"secret"`
}

// LoadAPIConfig 加载API配置
func LoadAPIConfig() (*APIConfig, error) {
	// 候选文件
	var realFile = Tea.ConfigFile(ConfigFileName)
	var oldRealFile = Tea.ConfigFile(oldConfigFileName)
	var isFromLocal = false
	var paths = []string{realFile, oldRealFile}
	homeDir, homeErr := os.UserHomeDir()
	if homeErr == nil {
		paths = append(paths, homeDir+"/."+teaconst.ProcessName+"/"+ConfigFileName)
	}
	paths = append(paths, "/etc/"+teaconst.ProcessName+"/"+ConfigFileName)

	var data []byte
	var err error
	var isFromOld = false
	for _, path := range paths {
		data, err = os.ReadFile(path)
		if err == nil {
			if path == realFile || path == oldRealFile {
				isFromLocal = true
			}

			// 自动生成新的配置文件
			isFromOld = path == oldRealFile

			break
		}
	}
	if err != nil {
		return nil, err
	}

	var config = &APIConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	err = config.Init()
	if err != nil {
		return nil, errors.New("init error: " + err.Error())
	}

	if !isFromLocal {
		// 恢复文件
		_ = os.WriteFile(realFile, data, 0666)
	}

	// 自动生成新配置文件
	if isFromOld {
		config.OldRPC.Endpoints = nil
		_ = config.WriteFile(Tea.ConfigFile(ConfigFileName))
	}

	return config, nil
}

// ResetAPIConfig 重置配置
func ResetAPIConfig() error {
	var filename = ConfigFileName

	// 重置 configs/api_admin.yaml
	{
		var configFile = Tea.ConfigFile(filename)
		stat, err := os.Stat(configFile)
		if err == nil && !stat.IsDir() {
			err = os.Remove(configFile)
			if err != nil {
				return err
			}
		}
	}

	// 重置 ~/.edge-admin/api_admin.yaml
	homeDir, homeErr := os.UserHomeDir()
	if homeErr == nil {
		var configFile = homeDir + "/." + teaconst.ProcessName + "/" + filename
		stat, err := os.Stat(configFile)
		if err == nil && !stat.IsDir() {
			err = os.Remove(configFile)
			if err != nil {
				return err
			}
		}
	}

	// 重置 /etc/edge-admin/api_admin.yaml
	{
		var configFile = "/etc/" + teaconst.ProcessName + "/" + filename
		stat, err := os.Stat(configFile)
		if err == nil && !stat.IsDir() {
			err = os.Remove(configFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func IsNewInstalled() bool {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	for _, filename := range []string{ConfigFileName, oldConfigFileName} {
		_, err = os.Stat(homeDir + "/." + teaconst.ProcessName + "/" + filename)
		if err == nil {
			return false
		}
	}

	return true
}

// WriteFile 写入API配置
func (this *APIConfig) WriteFile(path string) error {
	data, err := yaml.Marshal(this)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}

	// 写入 ~/ 和 /etc/ 目录，因为是备份需要，所以不需要提示错误信息
	// 写入 ~/.edge-admin/
	// 这个用来判断用户是否为重装，所以比较重要
	var filename = filepath.Base(path)
	homeDir, homeErr := os.UserHomeDir()
	if homeErr == nil {
		dir := homeDir + "/." + teaconst.ProcessName
		stat, err := os.Stat(dir)
		if err == nil && stat.IsDir() {
			err = os.WriteFile(dir+"/"+filename, data, 0666)
			if err != nil {
				return err
			}
		} else if err != nil && os.IsNotExist(err) {
			err = os.Mkdir(dir, 0777)
			if err == nil {
				err = os.WriteFile(dir+"/"+filename, data, 0666)
				if err != nil {
					return err
				}
			}
		}
	}

	// 写入 /etc/edge-admin
	{
		var dir = "/etc/" + teaconst.ProcessName
		stat, err := os.Stat(dir)
		if err == nil && stat.IsDir() {
			_ = os.WriteFile(dir+"/"+filename, data, 0666)
		} else if err != nil && os.IsNotExist(err) {
			err = os.Mkdir(dir, 0777)
			if err == nil {
				_ = os.WriteFile(dir+"/"+filename, data, 0666)
			}
		}
	}

	return nil
}

// Clone 克隆当前配置
func (this *APIConfig) Clone() *APIConfig {
	return &APIConfig{
		NodeId: this.NodeId,
		Secret: this.Secret,
	}
}

func (this *APIConfig) Init() error {
	// compatible with old
	if len(this.RPCEndpoints) == 0 && len(this.OldRPC.Endpoints) > 0 {
		this.RPCEndpoints = this.OldRPC.Endpoints
		this.RPCDisableUpdate = this.OldRPC.DisableUpdate
	}

	if len(this.RPCEndpoints) == 0 {
		return errors.New("no valid 'rpc.endpoints'")
	}

	if len(this.NodeId) == 0 {
		return errors.New("'nodeId' required")
	}
	if len(this.Secret) == 0 {
		return errors.New("'secret' required")
	}
	return nil
}
