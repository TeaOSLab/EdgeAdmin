package serverconfigs

import "encoding/json"

type ServerConfig struct {
	Id          string              `yaml:"id" json:"id"`                   // ID
	IsOn        bool                `yaml:"isOn" json:"isOn"`               // 是否开启
	Components  []*ComponentConfig  `yaml:"components" json:"components"`   // 组件
	Filters     []*FilterConfig     `yaml:"filters" json:"filters"`         // 过滤器
	Name        string              `yaml:"name" json:"name"`               // 名称
	Description string              `yaml:"description" json:"description"` // 描述
	ServerNames []*ServerNameConfig `yaml:"serverNames" json:"serverNames"` // 域名

	// 前端协议
	HTTP  *HTTPProtocolConfig  `yaml:"http" json:"http"`   // HTTP配置
	HTTPS *HTTPSProtocolConfig `yaml:"https" json:"https"` // HTTPS配置
	TCP   *TCPProtocolConfig   `yaml:"tcp" json:"tcp"`     // TCP配置
	TLS   *TLSProtocolConfig   `yaml:"tls" json:"tls"`     // TLS配置
	Unix  *UnixProtocolConfig  `yaml:"unix" json:"unix"`   // Unix配置
	UDP   *UDPProtocolConfig   `yaml:"udp" json:"udp"`     // UDP配置

	// Web配置
	Web *WebConfig `yaml:"web" json:"web"`

	// 反向代理配置
	ReverseProxy *ReverseProxyConfig `yaml:"reverseProxy" json:"reverseProxy"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{}
}

func (this *ServerConfig) Init() error {
	if this.HTTP != nil {
		err := this.HTTP.Init()
		if err != nil {
			return err
		}
	}

	if this.HTTPS != nil {
		err := this.HTTPS.Init()
		if err != nil {
			return err
		}
	}

	if this.TCP != nil {
		err := this.TCP.Init()
		if err != nil {
			return err
		}
	}

	if this.TLS != nil {
		err := this.TLS.Init()
		if err != nil {
			return err
		}
	}

	if this.Unix != nil {
		err := this.Unix.Init()
		if err != nil {
			return err
		}
	}

	if this.UDP != nil {
		err := this.UDP.Init()
		if err != nil {
			return err
		}
	}

	return nil
}

func (this *ServerConfig) FullAddresses() []string {
	result := []Protocol{}
	if this.HTTP != nil && this.HTTP.IsOn {
		result = append(result, this.HTTP.FullAddresses()...)
	}
	if this.HTTPS != nil && this.HTTPS.IsOn {
		result = append(result, this.HTTPS.FullAddresses()...)
	}
	if this.TCP != nil && this.TCP.IsOn {
		result = append(result, this.TCP.FullAddresses()...)
	}
	if this.TLS != nil && this.TLS.IsOn {
		result = append(result, this.TLS.FullAddresses()...)
	}
	if this.Unix != nil && this.Unix.IsOn {
		result = append(result, this.Unix.FullAddresses()...)
	}
	if this.UDP != nil && this.UDP.IsOn {
		result = append(result, this.UDP.FullAddresses()...)
	}

	return result
}

func (this *ServerConfig) Listen() []*NetworkAddressConfig {
	result := []*NetworkAddressConfig{}
	if this.HTTP != nil {
		result = append(result, this.HTTP.Listen...)
	}
	if this.HTTPS != nil {
		result = append(result, this.HTTPS.Listen...)
	}
	if this.TCP != nil {
		result = append(result, this.TCP.Listen...)
	}
	if this.TLS != nil {
		result = append(result, this.TLS.Listen...)
	}
	if this.Unix != nil {
		result = append(result, this.Unix.Listen...)
	}
	if this.UDP != nil {
		result = append(result, this.UDP.Listen...)
	}
	return result
}

func (this *ServerConfig) AsJSON() ([]byte, error) {
	return json.Marshal(this)
}

func (this *ServerConfig) IsHTTP() bool {
	return this.HTTP != nil || this.HTTPS != nil
}

func (this *ServerConfig) IsTCP() bool {
	return this.TCP != nil || this.TLS != nil
}

func (this *ServerConfig) IsUnix() bool {
	return this.Unix != nil
}

func (this *ServerConfig) IsUDP() bool {
	return this.UDP != nil
}
