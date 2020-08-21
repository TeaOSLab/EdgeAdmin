package serverconfigs

type ReverseProxyConfig struct {
	IsOn    bool                  `yaml:"isOn" json:"isOn"`       // 是否启用
	Origins []*OriginServerConfig `yaml:"origins" json:"origins"` // 源站列表
}
