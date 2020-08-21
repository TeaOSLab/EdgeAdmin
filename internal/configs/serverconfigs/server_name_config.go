package serverconfigs

type ServerNameType = string

const (
	ServerNameTypeFull   = "full"   // 完整的域名，包含通配符等
	ServerNameTypePrefix = "prefix" // 前缀
	ServerNameTypeSuffix = "suffix" // 后缀
	ServerNameTypeMatch  = "match"  // 正则匹配
)

type ServerNameConfig struct {
	Name string `yaml:"name" json:"name"` // 名称
	Type string `yaml:"type" json:"type"` // 类型
}
