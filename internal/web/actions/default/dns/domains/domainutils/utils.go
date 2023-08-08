package domainutils

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"net"
	"regexp"
	"strings"
)

// ValidateDomainFormat 校验域名格式
func ValidateDomainFormat(domain string) bool {
	pieces := strings.Split(domain, ".")
	for _, piece := range pieces {
		if piece == "-" ||
			strings.HasPrefix(piece, "-") ||
			strings.HasSuffix(piece, "-") ||
			len(piece) > 63 ||
			// 我们允许中文、大写字母、下划线，防止有些特殊场景下需要
			!regexp.MustCompile(`^[\p{Han}_a-zA-Z0-9-]+$`).MatchString(piece) {
			return false
		}
	}

	// 最后一段不能是全数字
	if regexp.MustCompile(`^(\d+)$`).MatchString(pieces[len(pieces)-1]) {
		return false
	}

	return true
}

// ConvertRoutesToMaps 转换线路列表
func ConvertRoutesToMaps(info *pb.NodeDNSInfo) []maps.Map {
	if info == nil {
		return []maps.Map{}
	}
	result := []maps.Map{}
	for _, route := range info.Routes {
		result = append(result, maps.Map{
			"name":       route.Name,
			"code":       route.Code,
			"domainId":   info.DnsDomainId,
			"domainName": info.DnsDomainName,
		})
	}
	return result
}

// FilterRoutes 筛选线路
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

// ValidateRecordName 校验记录名
func ValidateRecordName(name string) bool {
	if name == "*" || name == "@" || len(name) == 0 {
		return true
	}

	pieces := strings.Split(name, ".")
	for index, piece := range pieces {
		if index == 0 && piece == "*" {
			continue
		}
		if piece == "-" ||
			strings.HasPrefix(piece, "-") ||
			strings.HasSuffix(piece, "-") ||
			//strings.Contains(piece, "--") ||
			len(piece) > 63 ||
			// 我们允许中文、大写字母、下划线，防止有些特殊场景下需要
			!regexp.MustCompile(`^[\p{Han}_a-zA-Z0-9-]+$`).MatchString(piece) {
			return false
		}
	}
	return true
}

// ValidateRecordValue 校验记录值
func ValidateRecordValue(recordType dnsconfigs.RecordType, value string) (message string, ok bool) {
	switch recordType {
	case dnsconfigs.RecordTypeA:
		if !regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`).MatchString(value) {
			message = "请输入正确格式的IP"
			return
		}
		if net.ParseIP(value) == nil {
			message = "请输入正确格式的IP"
			return
		}
	case dnsconfigs.RecordTypeCNAME:
		value = strings.TrimSuffix(value, ".")
		if !strings.Contains(value, ".") || !ValidateDomainFormat(value) {
			message = "请输入正确的域名"
			return
		}
	case dnsconfigs.RecordTypeAAAA:
		if !strings.Contains(value, ":") {
			message = "请输入正确格式的IPv6地址"
			return
		}
		if net.ParseIP(value) == nil {
			message = "请输入正确格式的IPv6地址"
			return
		}
	case dnsconfigs.RecordTypeNS:
		value = strings.TrimSuffix(value, ".")
		if !strings.Contains(value, ".") || !ValidateDomainFormat(value) {
			message = "请输入正确的DNS服务器域名"
			return
		}
	case dnsconfigs.RecordTypeMX:
		value = strings.TrimSuffix(value, ".")
		if !strings.Contains(value, ".") || !ValidateDomainFormat(value) {
			message = "请输入正确的邮件服务器域名"
			return
		}
	case dnsconfigs.RecordTypeSRV:
		if len(value) == 0 {
			message = "请输入主机名"
			return
		}
	case dnsconfigs.RecordTypeTXT:
		if len(value) > 512 {
			message = "文本长度不能超出512字节"
			return
		}
	}

	if len(value) > 512 {
		message = "记录值长度不能超出512字节"
		return
	}

	ok = true
	return
}
