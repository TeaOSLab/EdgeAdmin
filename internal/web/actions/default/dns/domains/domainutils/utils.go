package domainutils

import (
	"regexp"
	"strings"
)

// 校验域名格式
func ValidateDomainFormat(domain string) bool {
	pieces := strings.Split(domain, ".")
	for _, piece := range pieces {
		if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(piece) {
			return false
		}
	}

	return true
}
