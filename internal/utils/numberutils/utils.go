package numberutils

import (
	"fmt"
	"github.com/iwind/TeaGo/types"
	"strconv"
)

func FormatInt64(value int64) string {
	return strconv.FormatInt(value, 10)
}

func FormatInt(value int) string {
	return strconv.Itoa(value)
}

func Pow1024(n int) int64 {
	if n <= 0 {
		return 1
	}
	if n == 1 {
		return 1024
	}
	return Pow1024(n-1) * 1024
}

func FormatBytes(bytes int64) string {
	if bytes < Pow1024(1) {
		return FormatInt64(bytes) + "B"
	} else if bytes < Pow1024(2) {
		return fmt.Sprintf("%.2fKB", float64(bytes)/float64(Pow1024(1)))
	} else if bytes < Pow1024(3) {
		return fmt.Sprintf("%.2fMB", float64(bytes)/float64(Pow1024(2)))
	} else if bytes < Pow1024(4) {
		return fmt.Sprintf("%.2fGB", float64(bytes)/float64(Pow1024(3)))
	} else if bytes < Pow1024(5) {
		return fmt.Sprintf("%.2fTB", float64(bytes)/float64(Pow1024(4)))
	} else if bytes < Pow1024(6) {
		return fmt.Sprintf("%.2fPB", float64(bytes)/float64(Pow1024(5)))
	} else {
		return fmt.Sprintf("%.2fEB", float64(bytes)/float64(Pow1024(6)))
	}
}

func FormatCount(count int64) string {
	if count < 1000 {
		return types.String(count)
	}
	if count < 1000*1000 {
		return fmt.Sprintf("%.1fK", float32(count)/1000)
	}
	if count < 1000*1000*1000 {
		return fmt.Sprintf("%.1fM", float32(count)/1000/1000)
	}
	return fmt.Sprintf("%.1fB", float32(count)/1000/1000/1000)
}

func FormatFloat(f interface{}, decimal int) string {
	if f == nil {
		return ""
	}
	switch x := f.(type) {
	case float32, float64:
		var s = fmt.Sprintf("%."+types.String(decimal)+"f", x)
		return s
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return types.String(x)
	case string:
		return x
	}
	return ""
}

func FormatFloat2(f interface{}) string {
	return FormatFloat(f, 2)
}
