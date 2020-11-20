package helpers

import (
	nodes "github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/securitymanager"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"net"
)

// 检查用户IP
func checkIP(config *securitymanager.SecurityConfig, ipAddr string) bool {
	if config == nil {
		return true
	}

	// 本地IP
	ip := net.ParseIP(ipAddr).To4()
	if ip == nil {
		logs.Println("[USER_MUST_AUTH]invalid client address: " + ipAddr)
		return false
	}
	if config.AllowLocal && isLocalIP(ip) {
		return true
	}

	// 检查位置
	if len(config.AllowCountryIds) > 0 || len(config.AllowProvinceIds) > 0 {
		rpc, err := nodes.SharedRPC()
		if err != nil {
			logs.Println("[USER_MUST_AUTH][ERROR]" + err.Error())
			return false
		}
		resp, err := rpc.IPLibraryRPC().LookupIPRegion(rpc.Context(0), &pb.LookupIPRegionRequest{Ip: ipAddr})
		if err != nil {
			logs.Println("[USER_MUST_AUTH][ERROR]" + err.Error())
			return false
		}
		if resp.Region == nil {
			return true
		}
		if len(config.AllowCountryIds) > 0 && !lists.ContainsInt64(config.AllowCountryIds, resp.Region.CountryId) {
			return false
		}
		if len(config.AllowProvinceIds) > 0 && !lists.ContainsInt64(config.AllowProvinceIds, resp.Region.ProvinceId) {
			return false
		}
	}

	return true
}

// 判断是否为本地IP
func isLocalIP(ip net.IP) bool {
	return ip[0] == 127 ||
		ip[0] == 10 ||
		(ip[0] == 172 && ip[1]&0xf0 == 16) ||
		(ip[0] == 192 && ip[1] == 168)
}
