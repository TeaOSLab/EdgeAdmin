package helpers

import (
	"github.com/iwind/TeaGo/actions"
	"net/http"
	"strconv"
)

type UserShouldAuth struct {
	action *actions.ActionObject
}

func (this *UserShouldAuth) BeforeAction(actionPtr actions.ActionWrapper, paramName string) (goNext bool) {
	this.action = actionPtr.Object()
	return true
}

// 存储用户名到SESSION
func (this *UserShouldAuth) StoreAdmin(adminId int, remember bool) {
	// 修改sid的时间
	if remember {
		cookie := &http.Cookie{
			Name:   "sid",
			Value:  this.action.Session().Sid,
			Path:   "/",
			MaxAge: 14 * 86400,
		}
		this.action.AddCookie(cookie)
	} else {
		cookie := &http.Cookie{
			Name:   "sid",
			Value:  this.action.Session().Sid,
			Path:   "/",
			MaxAge: 0,
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
