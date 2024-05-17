// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package group

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/groups/group/servergrouputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "group.index")
}

func (this *IndexAction) RunGet(params struct {
	GroupId int64
	Keyword string
}) {
	this.Data["keyword"] = params.Keyword

	group, err := servergrouputils.InitGroup(this.Parent(), params.GroupId, "")
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["group"] = maps.Map{
		"id":   group.Id,
		"name": group.Name,
	}

	// 是否有用户管理权限
	this.Data["canVisitUser"] = configloaders.AllowModule(this.AdminId(), configloaders.AdminModuleCodeUser)

	// 服务列表
	countResp, err := this.RPC().ServerRPC().CountAllEnabledServersMatch(this.AdminContext(), &pb.CountAllEnabledServersMatchRequest{
		ServerGroupId: params.GroupId,
		Keyword:       params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	var page = this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	// 服务列表
	serversResp, err := this.RPC().ServerRPC().ListEnabledServersMatch(this.AdminContext(), &pb.ListEnabledServersMatchRequest{
		Offset:         page.Offset,
		Size:           page.Size,
		ServerGroupId:  params.GroupId,
		Keyword:        params.Keyword,
		IgnoreSSLCerts: true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var serverMaps = []maps.Map{}
	for _, server := range serversResp.Servers {
		config := &serverconfigs.ServerConfig{}
		err = json.Unmarshal(server.Config, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 端口列表
		var portMaps = []maps.Map{}
		if config.HTTP != nil && config.HTTP.IsOn {
			for _, listen := range config.HTTP.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if config.HTTPS != nil && config.HTTPS.IsOn {
			for _, listen := range config.HTTPS.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if config.TCP != nil && config.TCP.IsOn {
			for _, listen := range config.TCP.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if config.TLS != nil && config.TLS.IsOn {
			for _, listen := range config.TLS.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if config.UDP != nil && config.UDP.IsOn {
			for _, listen := range config.UDP.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}

		// 分组
		var groupMaps = []maps.Map{}
		if len(server.ServerGroups) > 0 {
			for _, group := range server.ServerGroups {
				groupMaps = append(groupMaps, maps.Map{
					"id":   group.Id,
					"name": group.Name,
				})
			}
		}

		// 域名列表
		var serverNames = []*serverconfigs.ServerNameConfig{}
		if server.IsAuditing || (server.AuditingResult != nil && !server.AuditingResult.IsOk) {
			server.ServerNamesJSON = server.AuditingServerNamesJSON
		}
		var auditingIsOk = true
		if !server.IsAuditing && server.AuditingResult != nil && !server.AuditingResult.IsOk {
			auditingIsOk = false
		}
		if len(server.ServerNamesJSON) > 0 {
			err = json.Unmarshal(server.ServerNamesJSON, &serverNames)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
		var countServerNames = 0
		for _, serverName := range serverNames {
			if len(serverName.SubNames) == 0 {
				countServerNames++
			} else {
				countServerNames += len(serverName.SubNames)
			}
		}

		// 用户
		var userMap maps.Map = nil
		if server.User != nil {
			userMap = maps.Map{
				"id":       server.User.Id,
				"fullname": server.User.Fullname,
			}
		}

		serverMaps = append(serverMaps, maps.Map{
			"id":   server.Id,
			"isOn": server.IsOn,
			"name": server.Name,
			"cluster": maps.Map{
				"id":   server.NodeCluster.Id,
				"name": server.NodeCluster.Name,
			},
			"ports":            portMaps,
			"serverTypeName":   serverconfigs.FindServerType(server.Type).GetString("name"),
			"groups":           groupMaps,
			"serverNames":      serverNames,
			"countServerNames": countServerNames,
			"isAuditing":       server.IsAuditing,
			"auditingIsOk":     auditingIsOk,
			"user":             userMap,
		})
	}
	this.Data["servers"] = serverMaps

	this.Show()
}
