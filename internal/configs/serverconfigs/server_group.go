package serverconfigs

import "strings"

type ServerGroup struct {
	fullAddr string
	Servers  []*ServerConfig
}

func NewServerGroup(fullAddr string) *ServerGroup {
	return &ServerGroup{fullAddr: fullAddr}
}

// 添加服务
func (this *ServerGroup) Add(server *ServerConfig) {
	this.Servers = append(this.Servers, server)
}

// 获取完整的地址
func (this *ServerGroup) FullAddr() string {
	return this.fullAddr
}

// 获取当前分组的协议
func (this *ServerGroup) Protocol() Protocol {
	for _, p := range AllProtocols() {
		if strings.HasPrefix(this.fullAddr, p+":") {
			return p
		}
	}
	return ProtocolHTTP
}

// 获取当前分组的地址
func (this *ServerGroup) Addr() string {
	protocol := this.Protocol()
	if protocol == ProtocolUnix {
		return strings.TrimPrefix(this.fullAddr, protocol+":")
	}
	return strings.TrimPrefix(this.fullAddr, protocol+"://")
}
