package node

import (
	"fmt"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"strconv"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "node")
	this.SecondMenu("nodes")
}

func (this *IndexAction) RunGet(params struct {
	NodeId int64
}) {
	_, err := nodeutils.InitNodeInfo(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if teaconst.IsPlus {
		this.RedirectURL("/clusters/cluster/node/boards?clusterId=" + fmt.Sprintf("%d", this.Data["clusterId"]) + "&nodeId=" + strconv.FormatInt(params.NodeId, 10))
	} else {
		this.RedirectURL("/clusters/cluster/node/detail?clusterId=" + fmt.Sprintf("%d", this.Data["clusterId"]) + "&nodeId=" + strconv.FormatInt(params.NodeId, 10))
	}
}
