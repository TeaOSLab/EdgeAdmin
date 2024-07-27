package clusterutils

import (
	"net/http"

	"github.com/iwind/TeaGo/actions"
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
