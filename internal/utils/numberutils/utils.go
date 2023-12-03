package numberutils

import (
	"fmt"
	"github.com/iwind/TeaGo/types"
	"regexp"
	"strconv"
	"strings"
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
		return TrimZeroSuffix(fmt.Sprintf("%.2fKiB", float64(bytes)/float64(Pow1024(1))))
	} else if bytes < Pow1024(3) {
		return TrimZeroSuffix(fmt.Sprintf("%.2fMiB", float64(bytes)/float64(Pow1024(2))))
	} else if bytes < Pow1024(4) {
		return TrimZeroSuffix(fmt.Sprintf("%.2fGiB", float64(bytes)/float64(Pow1024(3))))
	} else if bytes < Pow1024(5) {
		return TrimZeroSuffix(fmt.Sprintf("%.2fTiB", float64(bytes)/float64(Pow1024(4))))
	} else if bytes < Pow1024(6) {
		return TrimZeroSuffix(fmt.Sprintf("%.2fPiB", float64(bytes)/float64(Pow1024(5))))
	} else {
		return TrimZeroSuffix(fmt.Sprintf("%.2fEiB", float64(bytes)/float64(Pow1024(6))))
	}
}

func FormatBits(bits int64) string {
	if bits < Pow1024(1) {
		return FormatInt64(bits) + "bps"
	} else if bits < Pow1024(2) {
		return TrimZeroSuffix(fmt.Sprintf("%.4fKbps", float64(bits)/float64(Pow1024(1))))
	} else if bits < Pow1024(3) {
		return TrimZeroSuffix(fmt.Sprintf("%.4fMbps", float64(bits)/float64(Pow1024(2))))
	} else if bits < Pow1024(4) {
		return TrimZeroSuffix(fmt.Sprintf("%.4fGbps", float64(bits)/float64(Pow1024(3))))
	} else if bits < Pow1024(5) {
		return TrimZeroSuffix(fmt.Sprintf("%.4fTbps", float64(bits)/float64(Pow1024(4))))
	} else if bits < Pow1024(6) {
		return TrimZeroSuffix(fmt.Sprintf("%.4fPbps", float64(bits)/float64(Pow1024(5))))
	} else {
		return TrimZeroSuffix(fmt.Sprintf("%.4fEbps", float64(bits)/float64(Pow1024(6))))
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

func FormatFloat(f any, decimal int) string {
	if f == nil {
		return ""
	}
	switch x := f.(type) {
	case float32, float64:
		var s = fmt.Sprintf("%."+types.String(decimal)+"f", x)

		// 分隔
		var dotIndex = strings.Index(s, ".")
		if dotIndex > 0 {
			var d = s[:dotIndex]
			var f2 = s[dotIndex:]
			f2 = strings.TrimRight(strings.TrimRight(f2, "0"), ".")
			return formatDigit(d) + f2
		}

		return s
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return formatDigit(types.String(x))
	case string:
		return x
	}
	return ""
}

func FormatFloat2(f any) string {
	return FormatFloat(f, 2)
}

// PadFloatZero 为浮点型数字字符串填充足够的0
func PadFloatZero(s string, countZero int) string {
	if countZero <= 0 {
		return s
	}
	if len(s) == 0 {
		s = "0"
	}
	var index = strings.Index(s, ".")
	if index < 0 {
		return s + "." + strings.Repeat("0", countZero)
	}
	var decimalLen = len(s) - 1 - index
	if decimalLen < countZero {
		return s + strings.Repeat("0", countZero-decimalLen)
	}
	return s
}

var decimalReg = regexp.MustCompile(`^(\d+\.\d+)([a-zA-Z]+)?$`)

// TrimZeroSuffix 去除小数数字尾部多余的0
func TrimZeroSuffix(s string) string {
	var matches = decimalReg.FindStringSubmatch(s)
	if len(matches) < 3 {
		return s
	}
	return strings.TrimRight(strings.TrimRight(matches[1], "0"), ".") + matches[2]
}

func formatDigit(d string) string {
	if len(d) == 0 {
		return d
	}

	var prefix = ""
	if d[0] < '0' || d[0] > '9' {
		prefix = d[:1]
		d = d[1:]
	}

	var l = len(d)
	if l > 3 {
		var pieces = l / 3
		var commIndex = l - pieces*3
		var d2 = ""
		if commIndex > 0 {
			d2 = d[:commIndex] + ", "
		}
		for i := 0; i < pieces; i++ {
			d2 += d[commIndex+i*3 : commIndex+i*3+3]
			if i != pieces-1 {
				d2 += ", "
			}
		}
		return prefix + d2
	}
	return prefix + d
}
