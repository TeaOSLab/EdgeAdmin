package ipadmin

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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
		"value":      item.Value,
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
	FirewallPolicyId int64
	ItemId           int64

	Value      string
	ExpiredAt  int64
	Reason     string
	Type       string
	EventLevel string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// 日志
	defer this.CreateLogInfo(codes.WAF_LogUpdateIPFromWAFPolicy, params.FirewallPolicyId, params.ItemId)

	switch params.Type {
	case "ip":
		// 校验IP格式
		params.Must.
			Field("value", params.Value).
			Require("请输入IP或IP段")

		_, _, _, ok := utils.ParseIPValue(params.Value)
		if !ok {
			this.FailField("value", "请输入正确的IP格式")
			return
		}
	case "all":
		params.Value = "0.0.0.0"
	}

	_, err := this.RPC().IPItemRPC().UpdateIPItem(this.AdminContext(), &pb.UpdateIPItemRequest{
		IpItemId:   params.ItemId,
		Value:      params.Value,
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
