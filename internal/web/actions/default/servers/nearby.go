// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"strings"
)

// NearbyAction 查找附近的Server
type NearbyAction struct {
	actionutils.ParentAction
}

func (this *NearbyAction) RunPost(params struct {
	ServerId int64
	Url      string
}) {
	var groupMaps = []maps.Map{}

	resp, err := this.RPC().ServerRPC().FindNearbyServers(this.AdminContext(), &pb.FindNearbyServersRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, group := range resp.Groups {
		switch resp.Scope {
		case "cluster":
			group.Name = "[集群]" + group.Name
		case "group":
			group.Name = "[分组]" + group.Name
		}

		var itemMaps = []maps.Map{}
		for _, server := range group.Servers {
			itemMaps = append(itemMaps, maps.Map{
				"name":     server.Name,
				"url":      strings.ReplaceAll(params.Url, "${serverId}", types.String(server.Id)),
				"isActive": params.ServerId == server.Id,
			})
		}

		groupMaps = append(groupMaps, maps.Map{
			"name":  group.Name,
			"items": itemMaps,
		})
	}

	this.Data["scope"] = resp.Scope
	this.Data["groups"] = groupMaps

	this.Success()
}
