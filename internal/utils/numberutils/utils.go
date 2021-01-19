package numberutils

import (
	"fmt"
	"strconv"
)

func FormatInt64(value int64) string {
	return strconv.FormatInt(value, 10)
}

func FormatInt(value int) string {
	return strconv.Itoa(value)
}

func FormatBytes(bytes int64) string {
	if bytes < 1024 {
		return FormatInt64(bytes) + "B"
	} else if bytes < 1024*1024 {
		return fmt.Sprintf("%.2fK", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2fM", float64(bytes)/1024/1024)
	} else if bytes < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fG", float64(bytes)/1024/1024/1024)
	} else {
		return fmt.Sprintf("%.2fP", float64(bytes)/1024/1024/1024/1024)
	}
}
