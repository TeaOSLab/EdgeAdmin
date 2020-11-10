package helpers

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/actions"
	"net/http"
	"strconv"
)

type UserShouldAuth struct {
	action *actions.ActionObject
}

func (this *UserShouldAuth) BeforeAction(actionPtr actions.ActionWrapper, paramName string) (goNext bool) {
	this.action = actionPtr.Object()

	// 安全相关
	action := this.action
	if !teaconst.EnabledFrame {
		action.AddHeader("X-Frame-Options", "SAMEORIGIN")
	}
	action.AddHeader("Content-Security-Policy", "default-src 'self' data:; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'")

	return true
}

// 存储用户名到SESSION
func (this *UserShouldAuth) StoreAdmin(adminId int, remember bool) {
	// 修改sid的时间
	if remember {
		cookie := &http.Cookie{
			Name:     "sid",
			Value:    this.action.Session().Sid,
			Path:     "/",
			MaxAge:   14 * 86400,
			HttpOnly: true,
		}
		if this.action.Request.TLS != nil {
			cookie.SameSite = http.SameSiteStrictMode
			cookie.Secure = true
		}
		this.action.AddCookie(cookie)
	} else {
		cookie := &http.Cookie{
			Name:     "sid",
			Value:    this.action.Session().Sid,
			Path:     "/",
			MaxAge:   0,
			HttpOnly: true,
		}
		if this.action.Request.TLS != nil {
			cookie.SameSite = http.SameSiteStrictMode
			cookie.Secure = true
		}
		this.action.AddCookie(cookie)
	}
	this.action.Session().Write("adminId", strconv.Itoa(adminId))
}

func (this *UserShouldAuth) IsUser() bool {
	return this.action.Session().GetInt("adminId") > 0
}

func (this *UserShouldAuth) AdminId() int {
	return this.action.Session().GetInt("adminId")
}

func (this *UserShouldAuth) Logout() {
	this.action.Session().Delete()
}
