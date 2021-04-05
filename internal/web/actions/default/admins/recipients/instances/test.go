package instances

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type TestAction struct {
	actionutils.ParentAction
}

func (this *TestAction) Init() {
	this.Nav("", "", "test")
}

func (this *TestAction) RunGet(params struct {
	InstanceId int64
}) {
	instanceResp, err := this.RPC().MessageMediaInstanceRPC().FindEnabledMessageMediaInstance(this.AdminContext(), &pb.FindEnabledMessageMediaInstanceRequest{MessageMediaInstanceId: params.InstanceId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	instance := instanceResp.MessageMediaInstance
	if instance == nil || instance.MessageMedia == nil {
		this.NotFound("messageMediaInstance", params.InstanceId)
		return
	}

	mediaParams := maps.Map{}
	if len(instance.ParamsJSON) > 0 {
		err = json.Unmarshal(instance.ParamsJSON, &mediaParams)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["instance"] = maps.Map{
		"id":   instance.Id,
		"isOn": instance.IsOn,
		"media": maps.Map{
			"type":            instance.MessageMedia.Type,
			"name":            instance.MessageMedia.Name,
			"userDescription": instance.MessageMedia.UserDescription,
		},
		"description": instance.Description,
		"params":      mediaParams,
	}

	this.Show()
}

func (this *TestAction) RunPost(params struct {
	InstanceId int64
	Subject    string
	Body       string
	User       string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("instanceId", params.InstanceId).
		Gt(0, "请选择正确的媒介")

	resp, err := this.RPC().MessageTaskRPC().CreateMessageTask(this.AdminContext(), &pb.CreateMessageTaskRequest{
		RecipientId: 0,
		InstanceId:  params.InstanceId,
		User:        params.User,
		Subject:     params.Subject,
		Body:        params.Body,
		IsPrimary:   true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["taskId"] = resp.MessageTaskId

	defer this.CreateLogInfo("创建媒介测试任务 %d", resp.MessageTaskId)

	this.Success()
}
