package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type CreateIPPopupAction struct {
	actionutils.ParentAction
}

func (this *CreateIPPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreateIPPopupAction) RunGet(params struct {
	ListId int64
	Type   string
}) {
	this.Data["listType"] = params.Type
	this.Data["listId"] = params.ListId

	this.Show()
}

func (this *CreateIPPopupAction) RunPost(params struct {
	ListId     int64
	IpFrom     string
	IpTo       string
	ExpiredAt  int64
	Reason     string
	Type       string
	EventLevel string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	switch params.Type {
	case "ipv4":
		params.Must.
			Field("ipFrom", params.IpFrom).
			Require("请输入开始IP")

		// 校验IP格式（ipFrom/ipTo）
		var ipFromLong uint64
		if !utils.IsIPv4(params.IpFrom) {
			this.Fail("请输入正确的开始IP")
		}
		ipFromLong = utils.IP2Long(params.IpFrom)

		var ipToLong uint64
		if len(params.IpTo) > 0 && !utils.IsIPv4(params.IpTo) {
			ipToLong = utils.IP2Long(params.IpTo)
			this.Fail("请输入正确的结束IP")
		}

		if ipFromLong > 0 && ipToLong > 0 && ipFromLong > ipToLong {
			params.IpTo, params.IpFrom = params.IpFrom, params.IpTo
		}
	case "ipv6":
		params.Must.
			Field("ipFrom", params.IpFrom).
			Require("请输入IP")

		// 校验IP格式（ipFrom）
		if !utils.IsIPv6(params.IpFrom) {
			this.Fail("请输入正确的IPv6地址")
		}
	case "all":
		params.IpFrom = "0.0.0.0"
	}

	createResp, err := this.RPC().IPItemRPC().CreateIPItem(this.AdminContext(), &pb.CreateIPItemRequest{
		IpListId:   params.ListId,
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
	itemId := createResp.IpItemId

	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "在WAF策略 %d 名单中添加IP %d", params.ListId, itemId)

	this.Success()
}
