package helpers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/events"
	nodes "github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"net"
	"sync"
)

var ipCacheMap = map[string]bool{} // ip => bool
var ipCacheLocker = sync.Mutex{}

func init() {
	events.On(events.EventSecurityConfigChanged, func() {
		ipCacheLocker.Lock()
		ipCacheMap = map[string]bool{}
		ipCacheLocker.Unlock()
	})
}

// 检查用户IP并支持缓存
func checkIP(config *systemconfigs.SecurityConfig, ipAddr string) bool {
	ipCacheLocker.Lock()
	ipCache, ok := ipCacheMap[ipAddr]
	if ok && ipCache {
		ipCacheLocker.Unlock()
		return ipCache
	}
	ipCacheLocker.Unlock()

	result := checkIPWithoutCache(config, ipAddr)
	ipCacheLocker.Lock()

	// 缓存的内容不能过多
	if len(ipCacheMap) > 100_000 {
		ipCacheMap = map[string]bool{}
	}

	ipCacheMap[ipAddr] = result
	ipCacheLocker.Unlock()
	return result
}

// 检查用户IP
func checkIPWithoutCache(config *systemconfigs.SecurityConfig, ipAddr string) bool {
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
