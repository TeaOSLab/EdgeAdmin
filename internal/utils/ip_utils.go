package utils

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/iwind/TeaGo/types"
	"math/big"
	"net"
	"regexp"
	"strings"
)

// IP2Long 将IP转换为整型
func IP2Long(ip string) uint64 {
	s := net.ParseIP(ip)
	if len(s) != 16 {
		return 0
	}

	if strings.Contains(ip, ":") { // IPv6
		bigInt := big.NewInt(0)
		bigInt.SetBytes(s.To16())
		return bigInt.Uint64()
	}
	return uint64(binary.BigEndian.Uint32(s.To4()))
}

// IsIPv4 判断是否为IPv4
func IsIPv4(ip string) bool {
	if !regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`).MatchString(ip) {
		return false
	}
	if IP2Long(ip) == 0 {
		return false
	}
	return true
}

// IsIPv6 判断是否为IPv6
func IsIPv6(ip string) bool {
	if !strings.Contains(ip, ":") {
		return false
	}
	return len(net.ParseIP(ip)) == net.IPv6len
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
