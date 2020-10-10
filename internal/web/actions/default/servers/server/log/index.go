package log

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "log", "")
	this.SecondMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	ServerId  int64
	RequestId string
}) {
	this.Data["serverId"] = params.ServerId
	this.Data["requestId"] = params.RequestId

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId  int64
	RequestId string

	Must *actions.Must
}) {
	isReverse := len(params.RequestId) > 0
	accessLogsResp, err := this.RPC().HTTPAccessLogRPC().ListHTTPAccessLogs(this.AdminContext(), &pb.ListHTTPAccessLogsRequest{
		ServerId:  params.ServerId,
		RequestId: params.RequestId,
		Size:      20,
		Day:       timeutil.Format("Ymd"),
		Reverse:   isReverse,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	accessLogs := accessLogsResp.AccessLogs
	if len(accessLogs) == 0 {
		accessLogs = []*pb.HTTPAccessLog{}
	}
	this.Data["accessLogs"] = accessLogs
	if len(accessLogs) > 0 {
		this.Data["requestId"] = accessLogs[0].RequestId
	} else {
		this.Data["requestId"] = params.RequestId
	}
	this.Data["hasMore"] = accessLogsResp.HasMore

	this.Success()
}
