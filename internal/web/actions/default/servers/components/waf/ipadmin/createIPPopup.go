package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/waf/ipadmin/ipadminutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/models"
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

	listId, err := models.SharedHTTPFirewallPolicyDAO.FindEnabledPolicyIPListIdWithType(this.AdminContext(), params.FirewallPolicyId, params.Type)
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
	// TODO 校验IP格式（ipFrom/ipTo）

	params.Must.
		Field("ipFrom", params.IpFrom).
		Require("请输入开始IP")

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
	this.CreateLog(oplogs.LevelInfo, "在WAF策略 %d 名单中添加IP %d", params.FirewallPolicyId, itemId)

	this.Success()
}
