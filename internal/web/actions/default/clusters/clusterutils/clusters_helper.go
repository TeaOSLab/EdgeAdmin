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

func (this *ClustersHelper) BeforeAction(actionPtr actions.ActionWrapper) {
	var action = actionPtr.Object()
	if action.Request.Method != http.MethodGet {
		return
	}

	action.Data["teaMenu"] = "clusters"
}
