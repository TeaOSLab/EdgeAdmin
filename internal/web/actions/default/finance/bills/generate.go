package bills

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type GenerateAction struct {
	actionutils.ParentAction
}

func (this *GenerateAction) Init() {
	this.Nav("", "", "generate")
}

func (this *GenerateAction) RunGet(params struct{}) {
	this.Data["month"] = timeutil.Format("Ym", time.Now().AddDate(0, -1, 0))

	this.Show()
}

func (this *GenerateAction) RunPost(params struct {
	Month string

	Must *actions.Must
}) {
	defer this.CreateLogInfo("手动生成上个月（" + params.Month + "）账单")

	_, err := this.RPC().UserBillRPC().GenerateAllUserBills(this.AdminContext(), &pb.GenerateAllUserBillsRequest{Month: params.Month})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
