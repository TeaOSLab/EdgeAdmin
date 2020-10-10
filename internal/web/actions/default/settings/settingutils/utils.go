package settingutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
)

type Helper struct {
	tab string
}

func NewHelper(tab string) *Helper {
	return &Helper{
		tab: tab,
	}
}

func (this *Helper) BeforeAction(actionPtr actions.ActionWrapper) (goNext bool) {
	goNext = true

	action := actionPtr.Object()

	// 左侧菜单
	action.Data["teaMenu"] = "settings"

	// 标签栏
	tabbar := actionutils.NewTabbar()
	tabbar.Add("管理界面", "", "/settings", "", this.tab == "console")
	tabbar.Add("安全设置", "", "/settings/security", "", this.tab == "security")
	tabbar.Add("API节点", "", "/api", "exchange", this.tab == "api")
	tabbar.Add("日志数据库节点", "", "/db", "database", this.tab == "db")
	tabbar.Add("备份", "", "/settings/backup", "", this.tab == "backup")
	tabbar.Add("个人资料", "", "/settings/profile", "", this.tab == "profile")
	tabbar.Add("登录设置", "", "/settings/login", "", this.tab == "login")
	tabbar.Add("检查版本更新", "", "/settings/upgrade", "", this.tab == "upgrade")
	actionutils.SetTabbar(actionPtr, tabbar)

	return
}
