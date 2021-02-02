package utils

import (
	"encoding/binary"
	"math/big"
	"net"
	"regexp"
	"strings"
)

// 将IP转换为整型
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

// 判断是否为IPv4
func IsIPv4(ip string) bool {
	if !regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`).MatchString(ip) {
		return false
	}
	if IP2Long(ip) == 0 {
		return false
	}
	return true
}

// 判断是否为IPv6
func IsIPv6(ip string) bool {
	if !strings.Contains(ip, ":") {
		return false
	}
	return len(net.ParseIP(ip)) == net.IPv6len
}
