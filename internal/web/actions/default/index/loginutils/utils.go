// Copyright 2023 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package loginutils

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/actions"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"net/http"
)

// CalculateClientFingerprint 计算客户端指纹
func CalculateClientFingerprint(action *actions.ActionObject) string {
	return stringutil.Md5(action.RequestRemoteIP() + "@" + action.Request.UserAgent())
}

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
