package recipients

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"regexp"
	"strings"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	RecipientId int64
}) {
	recipientResp, err := this.RPC().MessageRecipientRPC().FindEnabledMessageRecipient(this.AdminContext(), &pb.FindEnabledMessageRecipientRequest{MessageRecipientId: params.RecipientId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	recipient := recipientResp.MessageRecipient
	if recipient == nil || recipient.Admin == nil || recipient.MessageMediaInstance == nil {
		this.NotFound("messageRecipient", params.RecipientId)
		return
	}

	var timeFromHour = ""
	var timeFromMinute = ""
	var timeFromSecond = ""

	if len(recipient.TimeFrom) > 0 {
		var pieces = strings.Split(recipient.TimeFrom, ":")
		timeFromHour = pieces[0]
		timeFromMinute = pieces[1]
		timeFromSecond = pieces[2]
	}

	var timeToHour = ""
	var timeToMinute = ""
	var timeToSecond = ""
	if len(recipient.TimeTo) > 0 {
		var pieces = strings.Split(recipient.TimeTo, ":")
		timeToHour = pieces[0]
		timeToMinute = pieces[1]
		timeToSecond = pieces[2]
	}

	this.Data["recipient"] = maps.Map{
		"id": recipient.Id,
		"admin": maps.Map{
			"id":       recipient.Admin.Id,
			"fullname": recipient.Admin.Fullname,
			"username": recipient.Admin.Username,
		},
		"groups": recipient.MessageRecipientGroups,
		"isOn":   recipient.IsOn,
		"instance": maps.Map{
			"id":   recipient.MessageMediaInstance.Id,
			"name": recipient.MessageMediaInstance.Name,
		},
		"user":           recipient.User,
		"description":    recipient.Description,
		"timeFromHour":   timeFromHour,
		"timeFromMinute": timeFromMinute,
		"timeFromSecond": timeFromSecond,
		"timeToHour":     timeToHour,
		"timeToMinute":   timeToMinute,
		"timeToSecond":   timeToSecond,
	}

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	RecipientId int64

	AdminId    int64
	User       string
	InstanceId int64

	GroupIds    string
	Description string
	IsOn        bool

	TimeFromHour   string
	TimeFromMinute string
	TimeFromSecond string

	TimeToHour   string
	TimeToMinute string
	TimeToSecond string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改媒介接收人 %d", params.RecipientId)

	params.Must.
		Field("adminId", params.AdminId).
		Gt(0, "请选择系统用户").
		Field("instanceId", params.InstanceId).
		Gt(0, "请选择媒介")

	groupIds := utils.SplitNumbers(params.GroupIds)

	var digitReg = regexp.MustCompile(`^\d+$`)

	var timeFrom = ""
	if digitReg.MatchString(params.TimeFromHour) && digitReg.MatchString(params.TimeFromMinute) && digitReg.MatchString(params.TimeFromSecond) {
		timeFrom = params.TimeFromHour + ":" + params.TimeFromMinute + ":" + params.TimeFromSecond
	}

	var timeTo = ""
	if digitReg.MatchString(params.TimeToHour) && digitReg.MatchString(params.TimeToMinute) && digitReg.MatchString(params.TimeToSecond) {
		timeTo = params.TimeToHour + ":" + params.TimeToMinute + ":" + params.TimeToSecond
	}

	_, err := this.RPC().MessageRecipientRPC().UpdateMessageRecipient(this.AdminContext(), &pb.UpdateMessageRecipientRequest{
		MessageRecipientId:       params.RecipientId,
		AdminId:                  params.AdminId,
		MessageMediaInstanceId:   params.InstanceId,
		User:                     params.User,
		MessageRecipientGroupIds: groupIds,
		Description:              params.Description,
		IsOn:                     params.IsOn,
		TimeFrom:                 timeFrom,
		TimeTo:                   timeTo,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
