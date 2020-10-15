package utils

import "strings"

// format address
func FormatAddress(addr string) string {
	if strings.HasSuffix(addr, "unix:") {
		return addr
	}
	addr = strings.Replace(addr, " ", "", -1)
	addr = strings.Replace(addr, "\t", "", -1)
	addr = strings.Replace(addr, "：", ":", -1)
	addr = strings.TrimSpace(addr)
	return addr
}
