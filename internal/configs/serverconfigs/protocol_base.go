package serverconfigs

// 协议基础数据结构
type BaseProtocol struct {
	IsOn   bool                    `yaml:"isOn" json:"isOn"`     // 是否开启
	Listen []*NetworkAddressConfig `yaml:"listen" json:"listen"` // 绑定的网络地址
}

// 初始化
func (this *BaseProtocol) InitBase() error {
	for _, addr := range this.Listen {
		err := addr.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

// 获取完整的地址列表
func (this *BaseProtocol) FullAddresses() []string {
	result := []string{}
	for _, addr := range this.Listen {
		result = append(result, addr.FullAddresses()...)
	}
	return result
}

// 添加地址
func (this *BaseProtocol) AddListen(addr ...*NetworkAddressConfig) {
	this.Listen = append(this.Listen, addr...)
}
