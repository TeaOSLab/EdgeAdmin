package iplists

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net"
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

	IpFrom string
	IpTo   string

	IpData string

	ExpiredAt  int64
	Reason     string
	Type       string
	EventLevel string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// 校验IPList
	existsResp, err := this.RPC().IPListRPC().ExistsEnabledIPList(this.AdminContext(), &pb.ExistsEnabledIPListRequest{IpListId: params.ListId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if !existsResp.Exists {
		this.Fail("IP名单不存在")
	}

	type ipData struct {
		ipFrom string
		ipTo   string
	}

	var batchIPs = []*ipData{}
	switch params.Type {
	case "ipv4":
		if params.Method == "single" {
			// 校验IP格式（ipFrom/ipTo）
			params.Must.
				Field("ipFrom", params.IpFrom).
				Require("请输入开始IP")

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
		} else if params.Method == "batch" {
			if len(params.IpData) == 0 {
				this.FailField("ipData", "请输入IP")
			}
			var lines = strings.Split(params.IpData, "\n")
			for index, line := range lines {
				line = strings.TrimSpace(line)
				if strings.Contains(line, "/") { // CIDR
					if strings.Contains(line, ":") {
						this.FailField("ipData", "第"+types.String(index+1)+"行IP格式错误："+line)
					}
					ipFrom, ipTo, err := configutils.ParseCIDR(line)
					if err != nil {
						this.FailField("ipData", "第"+types.String(index+1)+"行IP格式错误："+line)
					}
					batchIPs = append(batchIPs, &ipData{
						ipFrom: ipFrom,
						ipTo:   ipTo,
					})
				} else if strings.Contains(line, "-") { // IP Range
					var pieces = strings.Split(line, "-")
					var ipFrom = strings.TrimSpace(pieces[0])
					var ipTo = strings.TrimSpace(pieces[1])

					if net.ParseIP(ipFrom) == nil || net.ParseIP(ipTo) == nil || strings.Contains(ipFrom, ":") || strings.Contains(ipTo, ":") {
						this.FailField("ipData", "第"+types.String(index+1)+"行IP格式错误："+line)
					}
					if utils.IP2Long(ipFrom) > utils.IP2Long(ipTo) {
						ipFrom, ipTo = ipTo, ipFrom
					}
					batchIPs = append(batchIPs, &ipData{
						ipFrom: ipFrom,
						ipTo:   ipTo,
					})
				} else if strings.Contains(line, ",") { // IP Range
					var pieces = strings.Split(line, ",")
					var ipFrom = strings.TrimSpace(pieces[0])
					var ipTo = strings.TrimSpace(pieces[1])

					if net.ParseIP(ipFrom) == nil || net.ParseIP(ipTo) == nil || strings.Contains(ipFrom, ":") || strings.Contains(ipTo, ":") {
						this.FailField("ipData", "第"+types.String(index+1)+"行IP格式错误："+line)
					}
					if utils.IP2Long(ipFrom) > utils.IP2Long(ipTo) {
						ipFrom, ipTo = ipTo, ipFrom
					}
					batchIPs = append(batchIPs, &ipData{
						ipFrom: ipFrom,
						ipTo:   ipTo,
					})
				} else if len(line) > 0 {
					var ipFrom = line
					if net.ParseIP(ipFrom) == nil || strings.Contains(ipFrom, ":") {
						this.FailField("ipData", "第"+types.String(index+1)+"行IP格式错误："+line)
					}
					batchIPs = append(batchIPs, &ipData{
						ipFrom: ipFrom,
					})
				}
			}
		}
	case "ipv6":
		if params.Method == "single" {
			params.Must.
				Field("ipFrom", params.IpFrom).
				Require("请输入IP")

			// 校验IP格式（ipFrom）
			if !utils.IsIPv6(params.IpFrom) {
				this.Fail("请输入正确的IPv6地址")
			}
		} else if params.Method == "batch" {
			if len(params.IpData) == 0 {
				this.FailField("ipData", "请输入IP")
			}
			var lines = strings.Split(params.IpData, "\n")
			for index, line := range lines {
				line = strings.TrimSpace(line)
				if strings.Contains(line, "/") { // CIDR
					if !strings.Contains(line, ":") {
						this.FailField("ipData", "第"+types.String(index+1)+"行IP格式错误："+line)
					}
					ipFrom, ipTo, err := configutils.ParseCIDR(line)
					if err != nil {
						this.FailField("ipData", "第"+types.String(index+1)+"行IP格式错误："+line)
					}
					batchIPs = append(batchIPs, &ipData{
						ipFrom: ipFrom,
						ipTo:   ipTo,
					})
				} else if strings.Contains(line, "-") { // IP Range
					var pieces = strings.Split(line, "-")
					var ipFrom = strings.TrimSpace(pieces[0])
					var ipTo = strings.TrimSpace(pieces[1])

					if net.ParseIP(ipFrom) == nil || net.ParseIP(ipTo) == nil || !strings.Contains(ipFrom, ":") || !strings.Contains(ipTo, ":") {
						this.FailField("ipData", "第"+types.String(index+1)+"行IP格式错误："+line)
					}
					if utils.IP2Long(ipFrom) > utils.IP2Long(ipTo) {
						ipFrom, ipTo = ipTo, ipFrom
					}
					batchIPs = append(batchIPs, &ipData{
						ipFrom: ipFrom,
						ipTo:   ipTo,
					})
				} else if strings.Contains(line, ",") { // IP Range
					var pieces = strings.Split(line, ",")
					var ipFrom = strings.TrimSpace(pieces[0])
					var ipTo = strings.TrimSpace(pieces[1])

					if net.ParseIP(ipFrom) == nil || net.ParseIP(ipTo) == nil || !strings.Contains(ipFrom, ":") || !strings.Contains(ipTo, ":") {
						this.FailField("ipData", "第"+types.String(index+1)+"行IP格式错误："+line)
					}
					if utils.IP2Long(ipFrom) > utils.IP2Long(ipTo) {
						ipFrom, ipTo = ipTo, ipFrom
					}
					batchIPs = append(batchIPs, &ipData{
						ipFrom: ipFrom,
						ipTo:   ipTo,
					})
				} else if len(line) > 0 {
					var ipFrom = line
					if net.ParseIP(ipFrom) == nil || !strings.Contains(ipFrom, ":") {
						this.FailField("ipData", "第"+types.String(index+1)+"行IP格式错误："+line)
					}
					batchIPs = append(batchIPs, &ipData{
						ipFrom: ipFrom,
					})
				}
			}
		}
	case "all":
		params.IpFrom = "0.0.0.0"
	}

	if len(batchIPs) > 0 {
		for _, ip := range batchIPs {
			_, err := this.RPC().IPItemRPC().CreateIPItem(this.AdminContext(), &pb.CreateIPItemRequest{
				IpListId:   params.ListId,
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
		defer this.CreateLog(oplogs.LevelInfo, "在IP名单中批量添加IP %d", params.ListId)
	} else {
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
		defer this.CreateLog(oplogs.LevelInfo, "在IP名单 %d 中添加IP %d", params.ListId, itemId)
	}

	this.Success()
}
