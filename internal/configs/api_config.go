package configs

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/Tea"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// APIConfig API配置
type APIConfig struct {
	RPC struct {
		Endpoints     []string `yaml:"endpoints"`
		DisableUpdate bool     `yaml:"disableUpdate"`
	} `yaml:"rpc"`
	NodeId string `yaml:"nodeId"`
	Secret string `yaml:"secret"`
}

// LoadAPIConfig 加载API配置
func LoadAPIConfig() (*APIConfig, error) {
	// 候选文件
	var localFile = Tea.ConfigFile("api.yaml")
	var isFromLocal = false
	var paths = []string{localFile}
	homeDir, homeErr := os.UserHomeDir()
	if homeErr == nil {
		paths = append(paths, homeDir+"/."+teaconst.ProcessName+"/api.yaml")
	}
	paths = append(paths, "/etc/"+teaconst.ProcessName+"/api.yaml")

	var data []byte
	var err error
	for _, path := range paths {
		data, err = os.ReadFile(path)
		if err == nil {
			if path == localFile {
				isFromLocal = true
			}
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

	if !isFromLocal {
		// 恢复文件
		_ = os.WriteFile(localFile, data, 0666)
	}

	return config, nil
}

// ResetAPIConfig 重置配置
func ResetAPIConfig() error {
	var filename = "api.yaml"

	// 重置 configs/api.yaml
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

	// 重置 ~/.edge-admin/api.yaml
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

	// 重置 /etc/edge-admin/api.yaml
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
		RPC: struct {
			Endpoints     []string `yaml:"endpoints"`
			DisableUpdate bool     `yaml:"disableUpdate"`
		}{},
		NodeId: this.NodeId,
		Secret: this.Secret,
	}
}
