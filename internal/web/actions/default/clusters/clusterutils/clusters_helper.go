package clusterutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type ClustersHelper struct {
}

func NewClustersHelper() *ClustersHelper {
	return &ClustersHelper{}
}

func (this *ClustersHelper) BeforeAction(action *actions.ActionObject) {
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaMenu"] = "clusters"

	selectedTabbar, _ := action.Data["mainTab"]

	tabbar := actionutils.NewTabbar()
	tabbar.Add("集群", "", "/clusters", "", selectedTabbar == "cluster")
	tabbar.Add("SSH认证", "", "/clusters/grants", "", selectedTabbar == "grant")
	actionutils.SetTabbar(action, tabbar)
}
