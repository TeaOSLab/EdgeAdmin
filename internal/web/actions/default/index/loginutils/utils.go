// Copyright 2023 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package loginutils

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
	"github.com/iwind/TeaGo/actions"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"net"
	"net/http"
)

// CalculateClientFingerprint 计算客户端指纹
func CalculateClientFingerprint(action *actions.ActionObject) string {
	return stringutil.Md5(RemoteIP(action) + "@" + action.Request.UserAgent())
}

// RemoteIP 获取客户端IP
// TODO 将来增加是否使用代理设置（即从X-Real-IP中获取IP）
func RemoteIP(action *actions.ActionObject) string {
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
