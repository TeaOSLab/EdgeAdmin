package recipients

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"regexp"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	AdminId    int64
	InstanceId int64
	User       string

	TelegramToken string

	GroupIds    string
	Description string

	TimeFromHour   string
	TimeFromMinute string
	TimeFromSecond string

	TimeToHour   string
	TimeToMinute string
	TimeToSecond string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
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

	resp, err := this.RPC().MessageRecipientRPC().CreateMessageRecipient(this.AdminContext(), &pb.CreateMessageRecipientRequest{
		AdminId:                  params.AdminId,
		MessageMediaInstanceId:   params.InstanceId,
		User:                     params.User,
		MessageRecipientGroupIds: groupIds,
		Description:              params.Description,
		TimeFrom:                 timeFrom,
		TimeTo:                   timeTo,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	defer this.CreateLogInfo("创建媒介接收人 %d", resp.MessageRecipientId)

	this.Success()
}
