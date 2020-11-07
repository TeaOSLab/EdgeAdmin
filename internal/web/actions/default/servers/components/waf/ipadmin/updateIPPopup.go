package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateIPPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateIPPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateIPPopupAction) RunGet(params struct {
	ItemId int64
}) {
	itemResp, err := this.RPC().IPItemRPC().FindEnabledIPItem(this.AdminContext(), &pb.FindEnabledIPItemRequest{IpItemId: params.ItemId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	item := itemResp.IpItem
	if item == nil {
		this.NotFound("ipItem", params.ItemId)
		return
	}

	this.Data["item"] = maps.Map{
		"id":        item.Id,
		"ipFrom":    item.IpFrom,
		"ipTo":      item.IpTo,
		"expiredAt": item.ExpiredAt,
		"reason":    item.Reason,
	}

	this.Show()
}

func (this *UpdateIPPopupAction) RunPost(params struct {
	ItemId int64

	IpFrom    string
	IpTo      string
	ExpiredAt int64
	Reason    string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// TODO 校验ItemId所属用户
	// TODO 校验IP格式（ipFrom/ipTo）

	params.Must.
		Field("ipFrom", params.IpFrom).
		Require("请输入开始IP")

	_, err := this.RPC().IPItemRPC().UpdateIPItem(this.AdminContext(), &pb.UpdateIPItemRequest{
		IpItemId:  params.ItemId,
		IpFrom:    params.IpFrom,
		IpTo:      params.IpTo,
		ExpiredAt: params.ExpiredAt,
		Reason:    params.Reason,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
