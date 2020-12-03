package dashboard

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	// 取得用户的权限
	module, ok := configloaders.FindFirstAdminModule(this.AdminId())
	if ok {
		for _, m := range configloaders.AllModuleMaps() {
			if m.GetString("code") == module {
				this.RedirectURL(m.GetString("url"))
				return
			}
		}
	}

	this.Show()
}
