package items

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/regions/regionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
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
	Name     string
	BitsFrom int64
	BitsTo   int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入名称").
		Field("bitsFrom", params.BitsFrom).
		Gte(0, "请输入不小于0的整数").
		Field("bitsTo", params.BitsTo).
		Gte(0, "请输入不小于0的整数")

	createResp, err := this.RPC().NodePriceItemRPC().CreateNodePriceItem(this.AdminContext(), &pb.CreateNodePriceItemRequest{
		Name:     params.Name,
		Type:     regionutils.PriceTypeTraffic,
		BitsFrom: params.BitsFrom * 1000 * 1000,
		BitsTo:   params.BitsTo * 1000 * 1000,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	defer this.CreateLogInfo("创建流量价格项目", createResp.NodePriceItemId)
	this.Success()
}
