package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/iputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type ClusterAction struct {
	actionutils.ParentAction
}

func (this *ClusterAction) Init() {
	this.Nav("", "", "")
}

func (this *ClusterAction) RunGet(params struct {
	ClusterId int64
}) {
	// 集群信息
	clusterResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeCluster(this.AdminContext(), &pb.FindEnabledNodeClusterRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var cluster = clusterResp.NodeCluster
	if cluster == nil {
		this.NotFound("nodeCluster", params.ClusterId)
		return
	}
	this.Data["cluster"] = maps.Map{
		"id":   cluster.Id,
		"name": cluster.Name,
	}

	// DNS信息
	dnsResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterDNS(this.AdminContext(), &pb.FindEnabledNodeClusterDNSRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var defaultRoute = dnsResp.DefaultRoute
	var domainName = ""
	var dnsMap = maps.Map{
		"dnsName":          dnsResp.Name,
		"domainId":         0,
		"domainName":       "",
		"providerId":       0,
		"providerName":     "",
		"providerTypeName": "",
	}
	if dnsResp.Domain != nil {
		domainName = dnsResp.Domain.Name
		dnsMap["domainId"] = dnsResp.Domain.Id
		dnsMap["domainName"] = dnsResp.Domain.Name
	}
	if dnsResp.Provider != nil {
		dnsMap["providerId"] = dnsResp.Provider.Id
		dnsMap["providerName"] = dnsResp.Provider.Name
		dnsMap["providerTypeName"] = dnsResp.Provider.TypeName
	}

	if len(dnsResp.CnameRecords) > 0 {
		dnsMap["cnameRecords"] = dnsResp.CnameRecords
	} else {
		dnsMap["cnameRecords"] = []string{}
	}

	this.Data["dnsInfo"] = dnsMap

	// 未安装的节点
	notInstalledNodesResp, err := this.RPC().NodeRPC().FindAllEnabledNodesDNSWithNodeClusterId(this.AdminContext(), &pb.FindAllEnabledNodesDNSWithNodeClusterIdRequest{
		NodeClusterId: params.ClusterId,
		IsInstalled:   false,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var allNodes = notInstalledNodesResp.Nodes

	// 节点DNS解析记录
	nodesResp, err := this.RPC().NodeRPC().FindAllEnabledNodesDNSWithNodeClusterId(this.AdminContext(), &pb.FindAllEnabledNodesDNSWithNodeClusterIdRequest{
		NodeClusterId: params.ClusterId,
		IsInstalled:   true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var installedNodeIdsMap = map[int64]bool{}
	for _, node := range nodesResp.Nodes {
		installedNodeIdsMap[node.Id] = true
	}

	allNodes = append(allNodes, nodesResp.Nodes...)

	var nodeMaps = []maps.Map{}
	for _, node := range allNodes {
		var isInstalled = installedNodeIdsMap[node.Id]

		if len(node.Routes) > 0 {
			for _, route := range node.Routes {
				// 检查是否已解析
				var isResolved = false
				if isInstalled && cluster.DnsDomainId > 0 && len(cluster.DnsName) > 0 && len(node.IpAddr) > 0 {
					var recordType = "A"
					if iputils.IsIPv6(node.IpAddr) {
						recordType = "AAAA"
					}
					checkResp, err := this.RPC().DNSDomainRPC().ExistDNSDomainRecord(this.AdminContext(), &pb.ExistDNSDomainRecordRequest{
						DnsDomainId: cluster.DnsDomainId,
						Name:        cluster.DnsName,
						Type:        recordType,
						Route:       route.Code,
						Value:       node.IpAddr,
					})
					if err != nil {
						this.ErrorPage(err)
						return
					}
					isResolved = checkResp.IsOk
				}

				nodeMaps = append(nodeMaps, maps.Map{
					"id":       node.Id,
					"name":     node.Name,
					"ipAddr":   node.IpAddr,
					"ipAddrId": node.NodeIPAddressId,
					"route": maps.Map{
						"name": route.Name,
						"code": route.Code,
					},
					"clusterId":   node.NodeClusterId,
					"isResolved":  isResolved,
					"isInstalled": isInstalled,
					"isBackup":    node.IsBackupForCluster || node.IsBackupForGroup,
					"isOffline":   node.IsOffline,
				})
			}
		} else {
			// 默认线路
			var isResolved = false
			if isInstalled && len(defaultRoute) > 0 {
				var recordType = "A"
				if iputils.IsIPv6(node.IpAddr) {
					recordType = "AAAA"
				}
				checkResp, err := this.RPC().DNSDomainRPC().ExistDNSDomainRecord(this.AdminContext(), &pb.ExistDNSDomainRecordRequest{
					DnsDomainId: cluster.DnsDomainId,
					Name:        cluster.DnsName,
					Type:        recordType,
					Route:       defaultRoute,
					Value:       node.IpAddr,
				})
				if err != nil {
					this.ErrorPage(err)
					return
				}
				isResolved = checkResp.IsOk
			}
			nodeMaps = append(nodeMaps, maps.Map{
				"id":       node.Id,
				"name":     node.Name,
				"ipAddr":   node.IpAddr,
				"ipAddrId": node.NodeIPAddressId,
				"route": maps.Map{
					"name": "",
					"code": "",
				},
				"clusterId":   node.NodeClusterId,
				"isResolved":  isResolved,
				"isInstalled": isInstalled,
				"isBackup":    node.IsBackupForCluster || node.IsBackupForGroup,
				"isOffline":   node.IsOffline,
			})
		}
	}
	this.Data["nodes"] = nodeMaps

	// 代理服务解析记录
	serversResp, err := this.RPC().ServerRPC().FindAllEnabledServersDNSWithNodeClusterId(this.AdminContext(), &pb.FindAllEnabledServersDNSWithNodeClusterIdRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var serverMaps = []maps.Map{}
	for _, server := range serversResp.Servers {
		// 检查是否已解析
		isResolved := false
		if cluster.DnsDomainId > 0 && len(cluster.DnsName) > 0 && len(server.DnsName) > 0 && len(domainName) > 0 {
			checkResp, err := this.RPC().DNSDomainRPC().ExistDNSDomainRecord(this.AdminContext(), &pb.ExistDNSDomainRecordRequest{
				DnsDomainId: cluster.DnsDomainId,
				Name:        server.DnsName,
				Type:        "CNAME",
				Value:       cluster.DnsName + "." + domainName,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			isResolved = checkResp.IsOk
		}

		serverMaps = append(serverMaps, maps.Map{
			"id":         server.Id,
			"name":       server.Name,
			"dnsName":    server.DnsName,
			"isResolved": isResolved,
		})
	}
	this.Data["servers"] = serverMaps

	// 检查解析记录是否有变化
	checkChangesResp, err := this.RPC().NodeClusterRPC().CheckNodeClusterDNSChanges(this.AdminContext(), &pb.CheckNodeClusterDNSChangesRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["dnsHasChanges"] = checkChangesResp.IsChanged

	// 需要解决的问题
	issuesResp, err := this.RPC().DNSRPC().FindAllDNSIssues(this.AdminContext(), &pb.FindAllDNSIssuesRequest{
		NodeClusterId: params.ClusterId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var issueMaps = []maps.Map{}
	for _, issue := range issuesResp.Issues {
		issueMaps = append(issueMaps, maps.Map{
			"target":      issue.Target,
			"targetId":    issue.TargetId,
			"type":        issue.Type,
			"description": issue.Description,
			"params":      issue.Params,
		})
	}
	this.Data["issues"] = issueMaps

	// 当前正在执行的任务
	resp, err := this.RPC().DNSTaskRPC().FindAllDoingDNSTasks(this.AdminContext(), &pb.FindAllDoingDNSTasksRequest{
		NodeClusterId: params.ClusterId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var taskMaps = []maps.Map{}
	for _, task := range resp.DnsTasks {
		var clusterMap maps.Map = nil
		var nodeMap maps.Map = nil
		var serverMap maps.Map = nil
		var domainMap maps.Map = nil

		if task.NodeCluster != nil {
			clusterMap = maps.Map{
				"id":   task.NodeCluster.Id,
				"name": task.NodeCluster.Name,
			}
		}
		if task.Node != nil {
			nodeMap = maps.Map{
				"id":   task.Node.Id,
				"name": task.Node.Name,
			}
		}
		if task.Server != nil {
			serverMap = maps.Map{
				"id":   task.Server.Id,
				"name": task.Server.Name,
			}
		}
		if task.DnsDomain != nil {
			domainMap = maps.Map{
				"id":   task.DnsDomain.Id,
				"name": task.DnsDomain.Name,
			}
		}

		taskMaps = append(taskMaps, maps.Map{
			"id":          task.Id,
			"type":        task.Type,
			"isDone":      task.IsDone,
			"isOk":        task.IsOk,
			"error":       task.Error,
			"updatedTime": timeutil.FormatTime("Y-m-d H:i:s", task.UpdatedAt),
			"cluster":     clusterMap,
			"node":        nodeMap,
			"server":      serverMap,
			"domain":      domainMap,
		})
	}
	this.Data["tasks"] = taskMaps

	this.Show()
}
