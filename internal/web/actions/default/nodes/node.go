package nodes

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/grants/grantutils"
	"github.com/iwind/TeaGo/maps"
)

type NodeAction struct {
	actionutils.ParentAction
}

func (this *NodeAction) Init() {
	this.Nav("", "node", "index")
}

func (this *NodeAction) RunGet(params struct {
	NodeId int64
}) {
	this.Data["nodeId"] = params.NodeId

	nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	node := nodeResp.Node
	if node == nil {
		this.WriteString("找不到要操作的节点")
		return
	}

	var clusterMap maps.Map = nil
	if node.Cluster != nil {
		clusterMap = maps.Map{
			"id":   node.Cluster.Id,
			"name": node.Cluster.Name,
		}
	}

	var loginMap maps.Map = nil
	if node.Login != nil {
		loginParams := maps.Map{}
		if len(node.Login.Params) > 0 {
			err = json.Unmarshal(node.Login.Params, &loginParams)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		grantMap := maps.Map{}
		grantId := loginParams.GetInt64("grantId")
		if grantId > 0 {
			grantResp, err := this.RPC().NodeGrantRPC().FindEnabledGrant(this.AdminContext(), &pb.FindEnabledGrantRequest{GrantId: grantId})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if grantResp.Grant != nil {
				grantMap = maps.Map{
					"id":         grantResp.Grant.Id,
					"name":       grantResp.Grant.Name,
					"method":     grantResp.Grant.Method,
					"methodName": grantutils.FindGrantMethodName(grantResp.Grant.Method),
				}
			}
		}

		loginMap = maps.Map{
			"id":     node.Login.Id,
			"name":   node.Login.Name,
			"type":   node.Login.Type,
			"params": loginParams,
			"grant":  grantMap,
		}
	}

	this.Data["node"] = maps.Map{
		"id":      node.Id,
		"name":    node.Name,
		"cluster": clusterMap,
		"login":   loginMap,
	}

	this.Show()
}
