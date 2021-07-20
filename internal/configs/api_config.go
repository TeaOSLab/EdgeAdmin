package configs

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/go-yaml/yaml"
	"github.com/iwind/TeaGo/Tea"
	"io/ioutil"
	"os"
	"path/filepath"
)

// APIConfig API配置
type APIConfig struct {
	RPC struct {
		Endpoints []string `yaml:"endpoints"`
	} `yaml:"rpc"`
	NodeId string `yaml:"nodeId"`
	Secret string `yaml:"secret"`
}

// LoadAPIConfig 加载API配置
func LoadAPIConfig() (*APIConfig, error) {
	// 候选文件
	localFile := Tea.ConfigFile("api.yaml")
	isFromLocal := false
	paths := []string{localFile}
	homeDir, homeErr := os.UserHomeDir()
	if homeErr == nil {
		paths = append(paths, homeDir+"/."+teaconst.ProcessName+"/api.yaml")
	}
	paths = append(paths, "/etc/"+teaconst.ProcessName+"/api.yaml")

	var data []byte
	var err error
	for _, path := range paths {
		data, err = ioutil.ReadFile(path)
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

	config := &APIConfig{}
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	if !isFromLocal {
		// 恢复文件
		_ = ioutil.WriteFile(localFile, data, 0666)
	}

	return config, nil
}

// ResetAPIConfig 重置配置
func ResetAPIConfig() error {
	filename := "api.yaml"

	{
		configFile := Tea.ConfigFile(filename)
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
		configFile := homeDir + "/." + teaconst.ProcessName + "/" + filename
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
		configFile := "/etc/" + teaconst.ProcessName + "/" + filename
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

	// 写入 ~/ 和 /etc/ 目录，因为是备份需要，所以不需要提示错误信息
	// 写入 ~/.edge-admin/
	filename := filepath.Base(path)
	homeDir, homeErr := os.UserHomeDir()
	if homeErr == nil {
		dir := homeDir + "/." + teaconst.ProcessName
		stat, err := os.Stat(dir)
		if err == nil && stat.IsDir() {
			_ = ioutil.WriteFile(dir+"/"+filename, data, 0666)
		} else if err != nil && os.IsNotExist(err) {
			err = os.Mkdir(dir, 0777)
			if err == nil {
				_ = ioutil.WriteFile(dir+"/"+filename, data, 0666)
			}
		}
	}

	// 写入 /etc/edge-admin
	{
		dir := "/etc/" + teaconst.ProcessName
		stat, err := os.Stat(dir)
		if err == nil && stat.IsDir() {
			_ = ioutil.WriteFile(dir+"/"+filename, data, 0666)
		} else if err != nil && os.IsNotExist(err) {
			err = os.Mkdir(dir, 0777)
			if err == nil {
				_ = ioutil.WriteFile(dir+"/"+filename, data, 0666)
			}
		}
	}

	err = ioutil.WriteFile(path, data, 0666)
	if err != nil {
		return err
	}

	return nil
}
