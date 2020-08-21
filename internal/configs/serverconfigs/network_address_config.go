package serverconfigs

import (
	"github.com/iwind/TeaGo/types"
	"regexp"
	"strconv"
	"strings"
)

var regexpSinglePort = regexp.MustCompile(`^\d+$`)

// 网络地址配置
type NetworkAddressConfig struct {
	Protocol  string `yaml:"protocol" json:"protocol"`   // 协议，http、tcp、tcp4、tcp6、unix、udp等
	Host      string `yaml:"host" json:"host"`           // 主机地址或主机名
	PortRange string `yaml:"portRange" json:"portRange"` // 端口范围，支持 8080、8080-8090、8080:8090

	minPort int
	maxPort int
}

func (this *NetworkAddressConfig) Init() error {
	// 8080
	if regexpSinglePort.MatchString(this.PortRange) {
		this.minPort = types.Int(this.PortRange)
		this.maxPort = this.minPort
		return nil
	}

	// 8080:8090
	if strings.Contains(this.PortRange, ":") {
		pieces := strings.SplitN(this.PortRange, ":", 2)
		minPort := types.Int(pieces[0])
		maxPort := types.Int(pieces[1])
		if minPort > maxPort {
			minPort, maxPort = maxPort, minPort
		}
		this.minPort = minPort
		this.maxPort = maxPort
		return nil
	}

	// 8080-8090
	if strings.Contains(this.PortRange, "-") {
		pieces := strings.SplitN(this.PortRange, "-", 2)
		minPort := types.Int(pieces[0])
		maxPort := types.Int(pieces[1])
		if minPort > maxPort {
			minPort, maxPort = maxPort, minPort
		}
		this.minPort = minPort
		this.maxPort = maxPort
		return nil
	}

	return nil
}

func (this *NetworkAddressConfig) FullAddresses() []string {
	if this.Protocol == ProtocolUnix {
		return []string{this.Protocol + ":" + this.Host}
	}

	result := []string{}
	for i := this.minPort; i <= this.maxPort; i++ {
		host := this.Host
		result = append(result, this.Protocol+"://"+host+":"+strconv.Itoa(i))
	}
	return result
}
