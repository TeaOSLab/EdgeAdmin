package serverconfigs

// TODO 需要实现
type OriginServerGroupConfig struct {
	Origins []*OriginServerConfig `yaml:"origins" json:"origins"` // 源站列表
}
