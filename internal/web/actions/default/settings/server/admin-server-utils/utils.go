package adminserverutils

import (
	"errors"
	"github.com/iwind/TeaGo"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/types"
	"gopkg.in/yaml.v3"
	"net"
	"os"
	"time"
)

var ServerConfigIsChanged = false

const configFilename = "server.yaml"

// LoadServerConfig 读取当前服务配置
func LoadServerConfig() (*TeaGo.ServerConfig, error) {
	configFile := Tea.ConfigFile(configFilename)
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var serverConfig = &TeaGo.ServerConfig{}
	err = yaml.Unmarshal(data, serverConfig)
	if err != nil {
		return nil, err
	}
	return serverConfig, nil
}

// WriteServerConfig 保存当前服务配置
func WriteServerConfig(serverConfig *TeaGo.ServerConfig) error {
	data, err := yaml.Marshal(serverConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(Tea.ConfigFile(configFilename), data, 0666)
	if err != nil {
		return err
	}

	ServerConfigIsChanged = true

	return nil
}

// ReadServerHTTPS 检查HTTPS地址
func ReadServerHTTPS() (port int, err error) {
	config, err := LoadServerConfig()
	if err != nil {
		return 0, err
	}
	if config == nil {
		return 0, errors.New("could not load server config")
	}

	if config.Https.On && len(config.Https.Listen) > 0 {
		for _, listen := range config.Https.Listen {
			_, portString, splitErr := net.SplitHostPort(listen)
			if splitErr == nil {
				var portInt = types.Int(portString)
				if portInt > 0 {
					// 是否已经启动
					checkErr := func() error {
						conn, connErr := net.DialTimeout("tcp", ":"+portString, 1*time.Second)
						if connErr != nil {
							return connErr
						}
						_ = conn.Close()
						return nil
					}()
					if checkErr != nil {
						continue
					}

					port = portInt
					err = nil
					break
				}
			}
		}
	}
	return
}
