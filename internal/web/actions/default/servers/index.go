package servers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "server", "index")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId    int64
	GroupId      int64
	Keyword      string
	AuditingFlag int32
	CheckDNS     bool
}) {
	this.Data["clusterId"] = params.ClusterId
	this.Data["groupId"] = params.GroupId
	this.Data["keyword"] = params.Keyword
	this.Data["auditingFlag"] = params.AuditingFlag
	this.Data["checkDNS"] = params.CheckDNS

	isSearching := params.AuditingFlag == 1 || params.ClusterId > 0 || params.GroupId > 0 || len(params.Keyword) > 0

	if params.AuditingFlag > 0 {
		this.Data["firstMenuItem"] = "auditing"
	}

	// 常用的服务
	latestServerMaps := []maps.Map{}
	if !isSearching {
		serversResp, err := this.RPC().ServerRPC().FindLatestServers(this.AdminContext(), &pb.FindLatestServersRequest{Size: 6})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, server := range serversResp.Servers {
			latestServerMaps = append(latestServerMaps, maps.Map{
				"id":   server.Id,
				"name": server.Name,
			})
		}
	}
	this.Data["latestServers"] = latestServerMaps

	// 审核中的数量
	countAuditingResp, err := this.RPC().ServerRPC().CountAllEnabledServersMatch(this.AdminContext(), &pb.CountAllEnabledServersMatchRequest{
		AuditingFlag: 1,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countAuditing"] = countAuditingResp.Count

	// 全部数量
	countResp, err := this.RPC().ServerRPC().CountAllEnabledServersMatch(this.AdminContext(), &pb.CountAllEnabledServersMatchRequest{
		NodeClusterId: params.ClusterId,
		ServerGroupId: params.GroupId,
		Keyword:       params.Keyword,
		AuditingFlag:  params.AuditingFlag,
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
		Offset:        page.Offset,
		Size:          page.Size,
		NodeClusterId: params.ClusterId,
		ServerGroupId: params.GroupId,
		Keyword:       params.Keyword,
		AuditingFlag:  params.AuditingFlag,
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
		if len(server.ServerGroups) > 0 {
			for _, group := range server.ServerGroups {
				groupMaps = append(groupMaps, maps.Map{
					"id":   group.Id,
					"name": group.Name,
				})
			}
		}

		// 域名列表
		serverNames := []*serverconfigs.ServerNameConfig{}
		if server.IsAuditing || (server.AuditingResult != nil && !server.AuditingResult.IsOk) {
			server.ServerNamesJSON = server.AuditingServerNamesJSON
		}
		auditingIsOk := true
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
		countServerNames := 0
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

	// 集群
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClusters(this.AdminContext(), &pb.FindAllEnabledNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range clustersResp.NodeClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
		})
	}
	this.Data["clusters"] = clusterMaps

	// 分组
	groupsResp, err := this.RPC().ServerGroupRPC().FindAllEnabledServerGroups(this.AdminContext(), &pb.FindAllEnabledServerGroupsRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	groupMaps := []maps.Map{}
	for _, group := range groupsResp.ServerGroups {
		groupName := group.Name
		groupMaps = append(groupMaps, maps.Map{
			"id":   group.Id,
			"name": groupName,
		})
	}
	this.Data["groups"] = groupMaps

	// 是否有用户管理权限
	this.Data["canVisitUser"] = configloaders.AllowModule(this.AdminId(), configloaders.AdminModuleCodeUser)

	// 显示服务相关的日志
	errorLogsResp, err := this.RPC().NodeLogRPC().ListNodeLogs(this.AdminContext(), &pb.ListNodeLogsRequest{
		NodeId:     0,
		Role:       nodeconfigs.NodeRoleNode,
		Offset:     0,
		Size:       10,
		Level:      "",
		FixedState: int32(configutils.BoolStateNo),
		AllServers: true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	errorLogMaps := []maps.Map{}
	for _, errorLog := range errorLogsResp.NodeLogs {
		serverResp, err := this.RPC().ServerRPC().FindEnabledUserServerBasic(this.AdminContext(), &pb.FindEnabledUserServerBasicRequest{ServerId: errorLog.ServerId})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 服务
		var server = serverResp.Server
		if server == nil {
			// 设置为已修复
			_, err = this.RPC().NodeLogRPC().FixNodeLog(this.AdminContext(), &pb.FixNodeLogRequest{NodeLogId: errorLog.Id})
			if err != nil {
				this.ErrorPage(err)
				return
			}

			continue
		}

		// 节点
		nodeResp, err := this.RPC().NodeRPC().FindEnabledNode(this.AdminContext(), &pb.FindEnabledNodeRequest{NodeId: errorLog.NodeId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var node = nodeResp.Node
		if node == nil || node.NodeCluster == nil {
			// 设置为已修复
			_, err = this.RPC().NodeLogRPC().FixNodeLog(this.AdminContext(), &pb.FixNodeLogRequest{NodeLogId: errorLog.Id})
			if err != nil {
				this.ErrorPage(err)
				return
			}

			continue
		}

		errorLogMaps = append(errorLogMaps, maps.Map{
			"id":          errorLog.Id,
			"description": errorLog.Description,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", errorLog.CreatedAt),
			"serverId":    errorLog.ServerId,
			"level":       errorLog.Level,
			"serverName":  server.Name,
			"nodeId":      node.Id,
			"nodeName":    node.Name,
			"clusterId":   node.NodeCluster.Id,
		})
	}
	this.Data["errorLogs"] = errorLogMaps

	this.Show()
}
