package serverconfigs

// 源站服务配置
type OriginServerConfig struct {
	Id          string                `yaml:"id" json:"id"`                   // ID
	IsOn        bool                  `yaml:"isOn" json:"isOn"`               // 是否启用
	Name        string                `yaml:"name" json:"name"`               // 名称 TODO
	Addr        *NetworkAddressConfig `yaml:"addr" json:"addr"`               // 地址
	Description string                `yaml:"description" json:"description"` // 描述 TODO
}
