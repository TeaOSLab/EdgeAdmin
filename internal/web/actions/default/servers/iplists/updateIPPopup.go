package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/iputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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
		"id":         item.Id,
		"ipFrom":     item.IpFrom,
		"ipTo":       item.IpTo,
		"expiredAt":  item.ExpiredAt,
		"reason":     item.Reason,
		"type":       item.Type,
		"eventLevel": item.EventLevel,
	}

	this.Data["type"] = item.Type

	this.Show()
}

func (this *UpdateIPPopupAction) RunPost(params struct {
	ItemId int64

	IpFrom     string
	IpTo       string
	ExpiredAt  int64
	Reason     string
	Type       string
	EventLevel string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// 日志
	defer this.CreateLogInfo(codes.IPItem_LogUpdateIPItem, params.ItemId)

	// TODO 校验ItemId所属用户

	switch params.Type {
	case "ipv4":
		params.Must.
			Field("ipFrom", params.IpFrom).
			Require("请输入开始IP")

		// 校验IP格式（ipFrom/ipTo）
		if !iputils.IsIPv4(params.IpFrom) {
			this.FailField("ipFrom", "请输入正确的开始IP")
		}

		if len(params.IpTo) > 0 && !iputils.IsIPv4(params.IpTo) {
			this.FailField("ipTo", "请输入正确的结束IP")
		}

		if len(params.IpTo) > 0 && iputils.CompareIP(params.IpFrom, params.IpTo) > 0 {
			params.IpTo, params.IpFrom = params.IpFrom, params.IpTo
		}
	case "ipv6":
		params.Must.
			Field("ipFrom", params.IpFrom).
			Require("请输入正确的开始IP")

		if !iputils.IsIPv6(params.IpFrom) {
			this.FailField("ipFrom", "请输入正确的IPv6地址")
		}

		if len(params.IpTo) > 0 && !iputils.IsIPv6(params.IpTo) {
			this.FailField("ipTo", "请输入正确的IPv6地址")
		}

		if len(params.IpTo) > 0 && iputils.CompareIP(params.IpFrom, params.IpTo) > 0 {
			params.IpTo, params.IpFrom = params.IpFrom, params.IpTo
		}
	case "all":
		params.IpFrom = "0.0.0.0"
	}

	_, err := this.RPC().IPItemRPC().UpdateIPItem(this.AdminContext(), &pb.UpdateIPItemRequest{
		IpItemId:   params.ItemId,
		IpFrom:     params.IpFrom,
		IpTo:       params.IpTo,
		ExpiredAt:  params.ExpiredAt,
		Reason:     params.Reason,
		Type:       params.Type,
		EventLevel: params.EventLevel,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
