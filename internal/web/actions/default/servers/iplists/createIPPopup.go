package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
)

type CreateIPPopupAction struct {
	actionutils.ParentAction
}

func (this *CreateIPPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreateIPPopupAction) RunGet(params struct {
	ListId int64
}) {
	this.Data["listId"] = params.ListId

	listResp, err := this.RPC().IPListRPC().FindEnabledIPList(this.AdminContext(), &pb.FindEnabledIPListRequest{
		IpListId: params.ListId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var ipList = listResp.IpList
	if ipList == nil {
		this.NotFound("ipList", params.ListId)
		return
	}
	this.Data["list"] = maps.Map{
		"type": ipList.Type,
	}

	this.Show()
}

func (this *CreateIPPopupAction) RunPost(params struct {
	ListId int64
	Method string

	Value  string
	IpData string

	ExpiredAt  int64
	Reason     string
	Type       string
	EventLevel string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// 校验IPList
	if !firewallconfigs.IsGlobalListId(params.ListId) {
		existsResp, err := this.RPC().IPListRPC().ExistsEnabledIPList(this.AdminContext(), &pb.ExistsEnabledIPListRequest{IpListId: params.ListId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if !existsResp.Exists {
			this.Fail("IP名单不存在")
		}
	}

	type ipData struct {
		value  string
		ipFrom string
		ipTo   string
	}

	var batchIPs = []*ipData{}
	switch params.Type {
	case "ip":
		if params.Method == "single" {
			// 校验IP格式
			params.Must.
				Field("value", params.Value).
				Require("请输入IP或IP段")

			_, _, _, ok := utils.ParseIPValue(params.Value)
			if !ok {
				this.FailField("value", "请输入正确的IP格式")
				return
			}
		} else if params.Method == "batch" {
			if len(params.IpData) == 0 {
				this.FailField("ipData", "请输入IP")
			}
			var lines = strings.Split(params.IpData, "\n")
			for index, line := range lines {
				line = strings.TrimSpace(line)
				if len(line) == 0 {
					continue
				}
				_, ipFrom, ipTo, ok := utils.ParseIPValue(line)
				if !ok {
					this.FailField("ipData", "第"+types.String(index+1)+"行IP格式错误："+line)
					return
				}
				batchIPs = append(batchIPs, &ipData{
					value:  line,
					ipFrom: ipFrom,
					ipTo:   ipTo,
				})
			}
		}
	case "all":
		params.Value = "0.0.0.0"
	}

	if len(batchIPs) > 0 {
		for _, ip := range batchIPs {
			_, err := this.RPC().IPItemRPC().CreateIPItem(this.AdminContext(), &pb.CreateIPItemRequest{
				IpListId:   params.ListId,
				Value:      ip.value,
				IpFrom:     ip.ipFrom,
				IpTo:       ip.ipTo,
				ExpiredAt:  params.ExpiredAt,
				Reason:     params.Reason,
				Type:       params.Type,
				EventLevel: params.EventLevel,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		// 日志
		defer this.CreateLogInfo(codes.IPList_LogCreateIPItemsBatch, params.ListId)
	} else {
		createResp, err := this.RPC().IPItemRPC().CreateIPItem(this.AdminContext(), &pb.CreateIPItemRequest{
			IpListId:   params.ListId,
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
		itemId := createResp.IpItemId

		// 日志
		defer this.CreateLogInfo(codes.IPItem_LogCreateIPItem, params.ListId, itemId)
	}

	this.Success()
}
