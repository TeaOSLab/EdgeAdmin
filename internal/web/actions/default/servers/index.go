package servers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
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
	UserId       int64

	TrafficOutOrder     string
	RequestsOrder       string
	AttackRequestsOrder string
}) {
	this.Data["clusterId"] = params.ClusterId
	this.Data["groupId"] = params.GroupId
	this.Data["keyword"] = params.Keyword
	this.Data["auditingFlag"] = params.AuditingFlag
	this.Data["checkDNS"] = params.CheckDNS
	this.Data["hasOrder"] = len(params.TrafficOutOrder) > 0
	this.Data["userId"] = params.UserId

	var isSearching = params.AuditingFlag == 1 || params.ClusterId > 0 || params.GroupId > 0 || len(params.Keyword) > 0

	if params.AuditingFlag > 0 {
		this.Data["firstMenuItem"] = "auditing"
	}

	// 常用的服务
	var latestServerMaps = []maps.Map{}
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
		UserId:        params.UserId,
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
		Offset:             page.Offset,
		Size:               page.Size,
		NodeClusterId:      params.ClusterId,
		ServerGroupId:      params.GroupId,
		Keyword:            params.Keyword,
		AuditingFlag:       params.AuditingFlag,
		TrafficOutDesc:     params.TrafficOutOrder == "desc",
		TrafficOutAsc:      params.TrafficOutOrder == "asc",
		RequestsAsc:        params.RequestsOrder == "asc",
		RequestsDesc:       params.RequestsOrder == "desc",
		AttackRequestsAsc:  params.AttackRequestsOrder == "asc",
		AttackRequestsDesc: params.AttackRequestsOrder == "desc",
		UserId:             params.UserId,
		IgnoreServerNames:  true,
		IgnoreSSLCerts:     true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var serverMaps = []maps.Map{}
	for _, server := range serversResp.Servers {
		var config = &serverconfigs.ServerConfig{}
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
		if server.IsAuditing || (server.AuditingResult != nil && !server.AuditingResult.IsOk) {
			server.ServerNamesJSON = server.AuditingServerNamesJSON

			if len(config.ServerNames) == 0 {
				// 审核中的域名
				if len(server.ServerNamesJSON) > 0 {
					var serverNames = []*serverconfigs.ServerNameConfig{}
					err = json.Unmarshal(server.ServerNamesJSON, &serverNames)
					if err != nil {
						this.ErrorPage(err)
						return
					}
					config.ServerNames = serverNames
				}
			}
		}
		var auditingIsOk = true
		if !server.IsAuditing && server.AuditingResult != nil && !server.AuditingResult.IsOk {
			auditingIsOk = false
		}
		var firstServerName = ""
		for _, serverNameConfig := range config.ServerNames {
			if len(serverNameConfig.Name) > 0 {
				firstServerName = serverNameConfig.Name
				break
			}
			if len(serverNameConfig.SubNames) > 0 {
				firstServerName = serverNameConfig.SubNames[0]
				break
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

		// 提交审核时间
		var auditingTime = ""
		if server.AuditingAt > 0 {
			auditingTime = timeutil.FormatTime("Y-m-d", server.AuditingAt)
		}

		// 统计数据
		var bandwidthBits int64 = 0
		if server.BandwidthBytes > 0 {
			bandwidthBits = server.BandwidthBytes * 8
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
			"firstServerName":  firstServerName,
			"countServerNames": server.CountServerNames,
			"isAuditing":       server.IsAuditing,
			"auditingIsOk":     auditingIsOk,
			"user":             userMap,
			"auditingTime":     auditingTime,
			"bandwidthBits":    bandwidthBits,
			"qps":              numberutils.FormatCount(server.CountRequests / 300),       /** 5 minutes **/
			"attackQPS":        numberutils.FormatCount(server.CountAttackRequests / 300), /** 5 minutes **/
		})
	}
	this.Data["servers"] = serverMaps

	// 集群
	clustersResp, err := this.RPC().NodeClusterRPC().FindAllEnabledNodeClusters(this.AdminContext(), &pb.FindAllEnabledNodeClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var clusterMaps = []maps.Map{}
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
	var groupMaps = []maps.Map{}
	for _, group := range groupsResp.ServerGroups {
		var groupName = group.Name
		groupMaps = append(groupMaps, maps.Map{
			"id":   group.Id,
			"name": groupName,
		})
	}
	this.Data["groups"] = groupMaps

	// 是否有用户管理权限
	this.Data["canVisitUser"] = configloaders.AllowModule(this.AdminId(), configloaders.AdminModuleCodeUser)

	// 显示服务需要修复的日志数量
	countNeedFixLogsResp, err := this.RPC().NodeLogRPC().CountNodeLogs(this.AdminContext(), &pb.CountNodeLogsRequest{
		Role:       nodeconfigs.NodeRoleNode,
		AllServers: true,
		FixedState: int32(configutils.BoolStateNo),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countNeedFixLogs"] = countNeedFixLogsResp.Count

	// 是否有用户
	countUsersResp, err := this.RPC().UserRPC().CountAllEnabledUsers(this.AdminContext(), &pb.CountAllEnabledUsersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasUsers"] = countUsersResp.Count > 0

	this.Show()
}
