package items

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	ItemId int64
}) {
	itemResp, err := this.RPC().NodePriceItemRPC().FindEnabledNodePriceItem(this.AdminContext(), &pb.FindEnabledNodePriceItemRequest{NodePriceItemId: params.ItemId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	item := itemResp.NodePriceItem
	if item == nil {
		this.NotFound("nodePriceItem", params.ItemId)
		return
	}

	this.Data["item"] = maps.Map{
		"id":       item.Id,
		"name":     item.Name,
		"bitsFrom": item.BitsFrom,
		"bitsTo":   item.BitsTo,
	}

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	ItemId   int64
	Name     string
	BitsFrom int64
	BitsTo   int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改流量价格项目", params.ItemId)

	params.Must.
		Field("name", params.Name).
		Require("请输入名称").
		Field("bitsFrom", params.BitsFrom).
		Gte(0, "请输入不小于0的整数").
		Field("bitsTo", params.BitsTo).
		Gte(0, "请输入不小于0的整数")

	_, err := this.RPC().NodePriceItemRPC().UpdateNodePriceItem(this.AdminContext(), &pb.UpdateNodePriceItemRequest{
		NodePriceItemId: params.ItemId,
		Name:            params.Name,
		BitsFrom:        params.BitsFrom * 1000 * 1000,
		BitsTo:          params.BitsTo * 1000 * 1000,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
