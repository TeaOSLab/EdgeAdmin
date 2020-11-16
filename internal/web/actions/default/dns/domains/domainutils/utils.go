package domainutils

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"regexp"
	"strings"
)

// 校验域名格式
func ValidateDomainFormat(domain string) bool {
	pieces := strings.Split(domain, ".")
	for _, piece := range pieces {
		if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(piece) {
			return false
		}
	}

	return true
}

// 转换线路列表
func ConvertRoutesToMaps(routes []*pb.DNSRoute) []maps.Map {
	result := []maps.Map{}
	for _, route := range routes {
		result = append(result, maps.Map{
			"name": route.Name,
			"code": route.Code,
		})
	}
	return result
}

// 筛选线路
func FilterRoutes(routes []*pb.DNSRoute, allRoutes []*pb.DNSRoute) []*pb.DNSRoute {
	routeCodes := []string{}
	for _, route := range allRoutes {
		routeCodes = append(routeCodes, route.Code)
	}
	result := []*pb.DNSRoute{}
	for _, route := range routes {
		if lists.ContainsString(routeCodes, route.Code) {
			result = append(result, route)
		}
	}
	return result
}
