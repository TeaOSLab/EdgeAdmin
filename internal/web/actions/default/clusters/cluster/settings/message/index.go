package message

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("message")
}

func (this *IndexAction) RunGet(params struct{}) {
	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId     int64
	ReceiversJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改集群 %d 消息接收人", params.ClusterId)

	receiverMaps := []maps.Map{}
	if len(params.ReceiversJSON) > 0 {
		err := json.Unmarshal(params.ReceiversJSON, &receiverMaps)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	pbReceiverOptions := &pb.UpdateMessageReceiversRequest_RecipientOptions{}
	for _, receiverMap := range receiverMaps {
		recipientId := int64(0)
		groupId := int64(0)
		receiverType := receiverMap.GetString("type")
		switch receiverType {
		case "recipient":
			recipientId = receiverMap.GetInt64("id")
		case "group":
			groupId = receiverMap.GetInt64("id")
		default:
			continue
		}
		pbReceiverOptions.RecipientOptions = append(pbReceiverOptions.RecipientOptions, &pb.UpdateMessageReceiversRequest_RecipientOption{
			MessageRecipientId:      recipientId,
			MessageRecipientGroupId: groupId,
		})
	}

	_, err := this.RPC().MessageReceiverRPC().UpdateMessageReceivers(this.AdminContext(), &pb.UpdateMessageReceiversRequest{
		NodeClusterId: params.ClusterId,
		NodeId:        0,
		ServerId:      0,
		ParamsJSON:    nil,
		RecipientOptions: map[string]*pb.UpdateMessageReceiversRequest_RecipientOptions{
			"*": pbReceiverOptions,
		},
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
