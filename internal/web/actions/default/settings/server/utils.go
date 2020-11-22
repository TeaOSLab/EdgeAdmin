package server

import (
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

var serverConfigIsChanged = false

// 读取当前服务配置
func loadServerConfig() (*TeaGo.ServerConfig, error) {
	configFile := Tea.ConfigFile("server.yaml")
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	serverConfig := &TeaGo.ServerConfig{}
	err = yaml.Unmarshal(data, serverConfig)
	if err != nil {
		return nil, err
	}
	return serverConfig, nil
}

// 保存当前服务配置
func writeServerConfig(serverConfig *TeaGo.ServerConfig) error {
	data, err := yaml.Marshal(serverConfig)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(Tea.ConfigFile("server.yaml"), data, 0666)
	if err != nil {
		return err
	}

	serverConfigIsChanged = true

	return nil
}
