package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/waf/ipadmin/ipadminutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
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
	FirewallPolicyId int64
	Type             string
}) {
	this.Data["type"] = params.Type

	listId, err := dao.SharedHTTPFirewallPolicyDAO.FindEnabledPolicyIPListIdWithType(this.AdminContext(), params.FirewallPolicyId, params.Type)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["listId"] = listId

	this.Show()
}

func (this *CreateIPPopupAction) RunPost(params struct {
	FirewallPolicyId int64
	ListId           int64
	IpFrom           string
	IpTo             string
	ExpiredAt        int64
	Reason           string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// TODO 校验ListId所属用户

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

	createResp, err := this.RPC().IPItemRPC().CreateIPItem(this.AdminContext(), &pb.CreateIPItemRequest{
		IpListId:  params.ListId,
		IpFrom:    params.IpFrom,
		IpTo:      params.IpTo,
		ExpiredAt: params.ExpiredAt,
		Reason:    params.Reason,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	itemId := createResp.IpItemId

	// 发送通知
	err = ipadminutils.NotifyUpdateToClustersWithFirewallPolicyId(this.AdminContext(), params.FirewallPolicyId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "在WAF策略 %d 名单中添加IP %d", params.FirewallPolicyId, itemId)

	this.Success()
}
