package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
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
	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "修改WAF策略名单中的IP %d", params.ItemId)

	// TODO 校验ItemId所属用户

	params.Must.
		Field("ipFrom", params.IpFrom).
		Require("请输入开始IP")

	// 校验IP格式（ipFrom/ipTo）
	ipFromLong := utils.IP2Long(params.IpFrom)
	if len(params.IpFrom) > 0 {
		if ipFromLong == 0 {
			this.Fail("请输入正确的开始IP")
		}
	}

	ipToLong := utils.IP2Long(params.IpTo)
	if len(params.IpTo) > 0 {
		if ipToLong == 0 {
			this.Fail("请输入正确的结束IP")
		}
	}

	if ipFromLong > 0 && ipToLong > 0 && ipFromLong > ipToLong {
		params.IpTo, params.IpFrom = params.IpFrom, params.IpTo
	}

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
