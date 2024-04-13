package utils

import (
	"bytes"
	"errors"
	"github.com/TeaOSLab/EdgeCommon/pkg/iputils"
	"github.com/iwind/TeaGo/types"
	"net"
	"strings"
)

// ParseIPValue 解析IP值
func ParseIPValue(value string) (newValue string, ipFrom string, ipTo string, ok bool) {
	if len(value) == 0 {
		return
	}

	newValue = value

	// ip1-ip2
	if strings.Contains(value, "-") {
		var pieces = strings.Split(value, "-")
		if len(pieces) != 2 {
			return
		}

		ipFrom = strings.TrimSpace(pieces[0])
		ipTo = strings.TrimSpace(pieces[1])

		if !iputils.IsValid(ipFrom) || !iputils.IsValid(ipTo) {
			return
		}

		if !iputils.IsSameVersion(ipFrom, ipTo) {
			return
		}

		if iputils.CompareIP(ipFrom, ipTo) > 0 {
			ipFrom, ipTo = ipTo, ipFrom
			newValue = ipFrom + "-" + ipTo
		}

		ok = true
		return
	}

	// ip/mask
	if strings.Contains(value, "/") {
		cidr, err := iputils.ParseCIDR(value)
		if err != nil {
			return
		}
		return newValue, cidr.From().String(), cidr.To().String(), true
	}

	// single value
	if iputils.IsValid(value) {
		ipFrom = value
		ok = true
		return
	}

	return
}

// ExtractIP 分解IP
// 只支持D段掩码的CIDR
// 最多只记录255个值
func ExtractIP(ipStrings string) ([]string, error) {
	ipStrings = strings.ReplaceAll(ipStrings, " ", "")

	// CIDR
	if strings.Contains(ipStrings, "/") {
		_, cidrNet, err := net.ParseCIDR(ipStrings)
		if err != nil {
			return nil, err
		}

		var index = strings.Index(ipStrings, "/")
		var ipFrom = ipStrings[:index]
		var bits = types.Int(ipStrings[index+1:])
		if bits < 24 {
			return nil, errors.New("CIDR bits should be greater than 24")
		}

		var ipv4 = net.ParseIP(ipFrom).To4()
		if len(ipv4) == 0 {
			return nil, errors.New("support IPv4 only")
		}

		var result = []string{}
		ipv4[3] = 0 // 从0开始
		for i := 0; i <= 255; i++ {
			if cidrNet.Contains(ipv4) {
				result = append(result, ipv4.String())
			}
			ipv4 = NextIP(ipv4)
		}
		return result, nil
	}

	// IP Range
	if strings.Contains(ipStrings, "-") {
		var index = strings.Index(ipStrings, "-")
		var ipFromString = ipStrings[:index]
		var ipToString = ipStrings[index+1:]

		var ipFrom = net.ParseIP(ipFromString).To4()
		if len(ipFrom) == 0 {
			return nil, errors.New("invalid ip '" + ipFromString + "'")
		}

		var ipTo = net.ParseIP(ipToString).To4()
		if len(ipTo) == 0 {
			return nil, errors.New("invalid ip '" + ipToString + "'")
		}

		if bytes.Compare(ipFrom, ipTo) > 0 {
			ipFrom, ipTo = ipTo, ipFrom
		}

		var result = []string{}
		for i := 0; i < 255; i++ {
			if bytes.Compare(ipFrom, ipTo) > 0 {
				break
			}
			result = append(result, ipFrom.String())
			ipFrom = NextIP(ipFrom)
		}
		return result, nil
	}

	return []string{ipStrings}, nil
}

// NextIP IP增加1
func NextIP(prevIP net.IP) net.IP {
	var ip = make(net.IP, len(prevIP))
	copy(ip, prevIP)
	var index = len(ip) - 1
	for {
		if ip[index] == 255 {
			ip[index] = 0
			index--
			if index < 0 {
				break
			}
		} else {
			ip[index]++
			break
		}
	}
	return ip
}

// IsLocalIP 判断是否为本地IP
// ip 是To4()或者To16()的结果
func IsLocalIP(ip net.IP) bool {
	if ip == nil {
		return false
	}

	if ip[0] == 127 ||
		ip[0] == 10 ||
		(ip[0] == 172 && ip[1]&0xf0 == 16) ||
		(ip[0] == 192 && ip[1] == 168) {
		return true
	}

	if ip.String() == "::1" {
		return true
	}
	return false
}
