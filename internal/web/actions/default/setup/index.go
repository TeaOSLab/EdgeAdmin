package setup

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"net"
	"regexp"
	"sort"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	// 当前服务器的IP
	serverIPs := []string{}
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		netAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		serverIPs = append(serverIPs, netAddr.IP.String())
	}

	// 对IP进行排序，我们希望IPv4排在前面，而且希望127.0.0.1排在IPv4里的最后
	sort.Slice(serverIPs, func(i, j int) bool {
		ip1 := serverIPs[i]

		if ip1 == "127.0.0.1" {
			return false
		}
		if regexp.MustCompile(`^\d+\.\d+\.\d+.\d+$`).MatchString(ip1) {
			return true
		}
		return false
	})
	this.Data["serverIPs"] = serverIPs

	this.Show()
}
