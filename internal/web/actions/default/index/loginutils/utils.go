// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package loginutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
	"github.com/iwind/TeaGo/actions"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"net"
	"net/http"
	"regexp"
	"strings"
)

// CalculateClientFingerprint 计算客户端指纹
func CalculateClientFingerprint(action *actions.ActionObject) string {
	return stringutil.Md5(RemoteIP(action) + "@" + action.Request.UserAgent())
}

// RemoteIP 获取客户端IP
func RemoteIP(action *actions.ActionObject) string {
	securityConfig, _ := configloaders.LoadSecurityConfig()

	if securityConfig != nil {
		if len(securityConfig.ClientIPHeaderNames) > 0 {
			var headerNames = regexp.MustCompile(`[,;\s，、；]`).Split(securityConfig.ClientIPHeaderNames, -1)
			for _, headerName := range headerNames {
				headerName = http.CanonicalHeaderKey(strings.TrimSpace(headerName))
				if len(headerName) == 0 {
					continue
				}

				var ipValue = action.Request.Header.Get(headerName)
				if net.ParseIP(ipValue) != nil {
					return ipValue
				}
			}

			if securityConfig.ClientIPHeaderOnly {
				return ""
			}
		}
	}

	ip, _, _ := net.SplitHostPort(action.Request.RemoteAddr)
	return ip
}

// LookupIPRegion 查找登录区域
func LookupIPRegion(ip string) string {
	if len(ip) == 0 {
		return ""
	}

	var result = iplibrary.LookupIP(ip)
	if result != nil && result.IsOk() {
		// 这里不需要网络运营商信息
		return result.CountryName() + "@" + result.ProvinceName() + "@" + result.CityName() + "@" + result.TownName()
	}

	return ""
}

// SetCookie 设置Cookie
func SetCookie(action *actions.ActionObject, remember bool) {
	if remember {
		var cookie = &http.Cookie{
			Name:     teaconst.CookieSID,
			Value:    action.Session().Sid,
			Path:     "/",
			MaxAge:   14 * 86400,
			HttpOnly: true,
		}
		if action.Request.TLS != nil {
			cookie.SameSite = http.SameSiteStrictMode
			cookie.Secure = true
		}
		action.AddCookie(cookie)
	} else {
		var cookie = &http.Cookie{
			Name:     teaconst.CookieSID,
			Value:    action.Session().Sid,
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		}
		if action.Request.TLS != nil {
			cookie.SameSite = http.SameSiteStrictMode
			cookie.Secure = true
		}
		action.AddCookie(cookie)
	}
}

// UnsetCookie 重置Cookie
func UnsetCookie(action *actions.ActionObject) {
	cookie := &http.Cookie{
		Name:     teaconst.CookieSID,
		Value:    action.Session().Sid,
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	if action.Request.TLS != nil {
		cookie.SameSite = http.SameSiteStrictMode
		cookie.Secure = true
	}
	action.AddCookie(cookie)
}
