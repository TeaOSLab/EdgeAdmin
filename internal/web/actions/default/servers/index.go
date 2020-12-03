package servers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	"strconv"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "server", "index")
}

func (this *IndexAction) RunGet(params struct {
	GroupId int64
	Keyword string
}) {
	this.Data["groupId"] = params.GroupId
	this.Data["keyword"] = params.Keyword

	countResp, err := this.RPC().ServerRPC().CountAllEnabledServersMatch(this.AdminContext(), &pb.CountAllEnabledServersMatchRequest{
		GroupId: params.GroupId,
		Keyword: params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	// 服务列表
	serversResp, err := this.RPC().ServerRPC().ListEnabledServersMatch(this.AdminContext(), &pb.ListEnabledServersMatchRequest{
		Offset:  page.Offset,
		Size:    page.Size,
		GroupId: params.GroupId,
		Keyword: params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	serverMaps := []maps.Map{}
	for _, server := range serversResp.Servers {
		config := &serverconfigs.ServerConfig{}
		err = json.Unmarshal(server.Config, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 端口列表
		portMaps := []maps.Map{}
		if len(server.HttpJSON) > 0 && config.HTTP.IsOn {
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
		if config.Unix != nil && config.Unix.IsOn {
			for _, listen := range config.Unix.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.Host,
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
		groupMaps := []maps.Map{}
		if len(server.Groups) > 0 {
			for _, group := range server.Groups {
				groupMaps = append(groupMaps, maps.Map{
					"id":   group.Id,
					"name": group.Name,
				})
			}
		}

		// 域名列表
		serverNames := []*serverconfigs.ServerNameConfig{}
		if len(server.ServerNamesJSON) > 0 {
			err = json.Unmarshal(server.ServerNamesJSON, &serverNames)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
		countServerNames := 0
		for _, serverName := range serverNames {
			if len(serverName.SubNames) == 0 {
				countServerNames++
			} else {
				countServerNames += len(serverName.SubNames)
			}
		}

		serverMaps = append(serverMaps, maps.Map{
			"id":   server.Id,
			"isOn": server.IsOn,
			"name": server.Name,
			"cluster": maps.Map{
				"id":   server.Cluster.Id,
				"name": server.Cluster.Name,
			},
			"ports":            portMaps,
			"serverTypeName":   serverconfigs.FindServerType(server.Type).GetString("name"),
			"groups":           groupMaps,
			"serverNames":      serverNames,
			"countServerNames": countServerNames,
		})
	}
	this.Data["servers"] = serverMaps

	// 分组
	groupsResp, err := this.RPC().ServerGroupRPC().FindAllEnabledServerGroups(this.AdminContext(), &pb.FindAllEnabledServerGroupsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	groupMaps := []maps.Map{}
	for _, group := range groupsResp.Groups {
		countResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithGroupId(this.AdminContext(), &pb.CountAllEnabledServersWithGroupIdRequest{GroupId: group.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		groupName := group.Name
		if countResp.Count > 0 {
			groupName += "(" + strconv.FormatInt(countResp.Count, 10) + ")"
		}
		groupMaps = append(groupMaps, maps.Map{
			"id":   group.Id,
			"name": groupName,
		})
	}
	this.Data["groups"] = groupMaps

	this.Show()
}
