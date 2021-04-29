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
		return fmt.Sprintf("%.2fKB", float64(bytes)/1024)
	} else if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.2fMB", float64(bytes)/1024/1024)
	} else if bytes < 1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fGB", float64(bytes)/1024/1024/1024)
	} else if bytes < 1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fTB", float64(bytes)/1024/1024/1024/1024)
	} else if bytes < 1024*1024*1024*1024*1024*1024 {
		return fmt.Sprintf("%.2fPB", float64(bytes)/1024/1024/1024/1024/1024)
	} else {
		return fmt.Sprintf("%.2fEB", float64(bytes)/1024/1024/1024/1024/1024/1024)
	}
}

func FormatBits(bits int64) string {
	if bits < 1000 {
		return FormatInt64(bits) + "B"
	} else if bits < 1000*1000 {
		return fmt.Sprintf("%.2fKB", float64(bits)/1000)
	} else if bits < 1000*1000*1000 {
		return fmt.Sprintf("%.2fMB", float64(bits)/1000/1000)
	} else if bits < 1000*1000*1000*1000 {
		return fmt.Sprintf("%.2fGB", float64(bits)/1000/1000/1000)
	} else if bits < 1000*1000*1000*1000*1000 {
		return fmt.Sprintf("%.2fTB", float64(bits)/1000/1000/1000/1000)
	} else if bits < 1000*1000*1000*1000*1000*1000 {
		return fmt.Sprintf("%.2fPB", float64(bits)/1000/1000/1000/1000/1000)
	} else {
		return fmt.Sprintf("%.2fEB", float64(bits)/1000/1000/1000/1000/1000/1000)
	}
}
