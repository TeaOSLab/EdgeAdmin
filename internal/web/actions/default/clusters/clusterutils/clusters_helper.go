package clusterutils

import (
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
}
