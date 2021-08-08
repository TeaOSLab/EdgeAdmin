package messages

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().MessageRPC().CountUnreadMessages(this.AdminContext(), &pb.CountUnreadMessagesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count

	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	listResp, err := this.RPC().MessageRPC().ListUnreadMessages(this.AdminContext(), &pb.ListUnreadMessagesRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	messages := []maps.Map{}
	for _, message := range listResp.Messages {
		clusterMap := maps.Map{}
		if message.NodeCluster != nil {
			clusterMap = maps.Map{
				"id":   message.NodeCluster.Id,
				"name": message.NodeCluster.Name,
			}
		}

		nodeMap := maps.Map{}
		if message.Node != nil {
			nodeMap = maps.Map{
				"id":   message.Node.Id,
				"name": message.Node.Name,
			}
		}

		messages = append(messages, maps.Map{
			"id":       message.Id,
			"role":     message.Role,
			"isRead":   message.IsRead,
			"body":     message.Body,
			"level":    message.Level,
			"datetime": timeutil.FormatTime("Y-m-d H:i:s", message.CreatedAt),
			"params":   string(message.ParamsJSON),
			"type":     message.Type,
			"cluster":  clusterMap,
			"node":     nodeMap,
		})
	}
	this.Data["messages"] = messages

	this.Show()
}
