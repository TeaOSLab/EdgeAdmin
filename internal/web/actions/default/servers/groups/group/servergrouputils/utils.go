// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package servergrouputils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
)

// InitGroup 初始化分组信息
func InitGroup(parent *actionutils.ParentAction, groupId int64, menuItem string) (*pb.ServerGroup, error) {
	groupResp, err := parent.RPC().ServerGroupRPC().FindEnabledServerGroup(parent.AdminContext(), &pb.FindEnabledServerGroupRequest{ServerGroupId: groupId})
	if err != nil {
		return nil, err
	}
	var group = groupResp.ServerGroup
	if group == nil {
		return nil, errors.New("group with id '" + types.String(groupId) + "' not found")
	}

	parent.Data["group"] = maps.Map{
		"id":   group.Id,
		"name": group.Name,
	}

	// 初始化设置菜单
	if len(menuItem) > 0 {
		// 获取设置概要信息
		configInfoResp, err := parent.RPC().ServerGroupRPC().FindEnabledServerGroupConfigInfo(parent.AdminContext(), &pb.FindEnabledServerGroupConfigInfoRequest{ServerGroupId: groupId})
		if err != nil {
			return group, err
		}

		var urlPrefix = "/servers/groups/group/settings"
		parent.Data["leftMenuItems"] = []maps.Map{
			/**{
				"name":     "Web设置",
				"url":      urlPrefix + "/web?groupId=" + types.String(groupId),
				"isActive": menuItem == "web",
			},**/
			{
				"name":     "HTTP反向代理",
				"url":      urlPrefix + "/httpReverseProxy?groupId=" + types.String(groupId),
				"isActive": menuItem == "httpReverseProxy",
				"isOn":     configInfoResp.HasHTTPReverseProxy,
			},
			{
				"name":     "TCP反向代理",
				"url":      urlPrefix + "/tcpReverseProxy?groupId=" + types.String(groupId),
				"isActive": menuItem == "tcpReverseProxy",
				"isOn":     configInfoResp.HasTCPReverseProxy,
			},
			{
				"name":     "UDP反向代理",
				"url":      urlPrefix + "/udpReverseProxy?groupId=" + types.String(groupId),
				"isActive": menuItem == "udpReverseProxy",
				"isOn":     configInfoResp.HasUDPReverseProxy,
			},
			/**{
				"name": "-",
				"url":  "",
			},
			{
				"name":     "WAF",
				"url":      urlPrefix + "/waf?groupId=" + types.String(groupId),
				"isActive": menuItem == "waf",
			},
			{
				"name":     "缓存",
				"url":      urlPrefix + "/cache?groupId=" + types.String(groupId),
				"isActive": menuItem == "cache",
			},
			{
				"name":     "访问日志",
				"url":      urlPrefix + "/accessLog?groupId=" + types.String(groupId),
				"isActive": menuItem == "accessLog",
			},
			{
				"name":     "统计",
				"url":      urlPrefix + "/stat?groupId=" + types.String(groupId),
				"isActive": menuItem == "stat",
			},
			{
				"name":     "Gzip压缩",
				"url":      urlPrefix + "/gzip?groupId=" + types.String(groupId),
				"isActive": menuItem == "gzip",
			},
			{
				"name":     "特殊页面",
				"url":      urlPrefix + "/pages?groupId=" + types.String(groupId),
				"isActive": menuItem == "page",
			},
			{
				"name":     "HTTP Header",
				"url":      urlPrefix + "/headers?groupId=" + types.String(groupId),
				"isActive": menuItem == "header",
			},
			{
				"name":     "Websocket",
				"url":      urlPrefix + "/websocket?groupId=" + types.String(groupId),
				"isActive": menuItem == "websocket",
			},**/
		}
	}

	return group, nil
}
