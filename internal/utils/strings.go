package utils

import (
	"github.com/iwind/TeaGo/types"
	"strings"
)

// FormatAddress format address
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

// SplitNumbers 分割数字
func SplitNumbers(numbers string) (result []int64) {
	if len(numbers) == 0 {
		return
	}
	pieces := strings.Split(numbers, ",")
	for _, piece := range pieces {
		number := types.Int64(strings.TrimSpace(piece))
		result = append(result, number)
	}
	return
}
