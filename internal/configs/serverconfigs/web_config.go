package serverconfigs

type WebConfig struct {
	IsOn bool `yaml:"isOn" json:"isOn"`

	Locations []*LocationConfig `yaml:"locations" json:"locations"` // 路径规则 TODO

	// 本地静态资源配置
	Root string `yaml:"root" json:"root"` // 资源根目录 TODO
}
