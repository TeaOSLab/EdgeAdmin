package utils

import (
	"encoding/binary"
	"net"
)

// 将IP转换为整型
func IP2Long(ip string) uint32 {
	s := net.ParseIP(ip)
	if s == nil {
		return 0
	}

	if len(s) == 16 {
		return binary.BigEndian.Uint32(s[12:16])
	}
	return binary.BigEndian.Uint32(s)
}
