// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package transfer

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/messageconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/types"
)

type UpgradeNodesAction struct {
	actionutils.ParentAction
}

func (this *UpgradeNodesAction) RunPost(params struct {
	ApiNodeProtocol string
	ApiNodeHost     string
	ApiNodePort     int
}) {
	nodesResp, err := this.RPC().NodeRPC().ListEnabledNodesMatch(this.AdminContext(), &pb.ListEnabledNodesMatchRequest{
		ActiveState: 1,
		Size:        100,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var nodes = nodesResp.Nodes
	this.Data["hasNext"] = len(nodes) > 0
	this.Data["count"] = len(nodes)

	if len(nodes) > 0 {
		var message = &messageconfigs.ChangeAPINodeMessage{
			Addr: params.ApiNodeProtocol + "://" + configutils.QuoteIP(params.ApiNodeHost) + ":" + types.String(params.ApiNodePort),
		}
		messageJSON, err := json.Marshal(message)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		for _, node := range nodesResp.Nodes {
			resp, err := this.RPC().NodeRPC().SendCommandToNode(this.AdminContext(), &pb.NodeStreamMessage{
				NodeId:         node.Id,
				TimeoutSeconds: 3,
				Code:           messageconfigs.MessageCodeChangeAPINode,
				DataJSON:       messageJSON,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if !resp.IsOk {
				this.Fail(resp.Message)
				return
			}
		}
	}

	this.Success()
}
