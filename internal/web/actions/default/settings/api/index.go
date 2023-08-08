package api

import (
	"context"
	"encoding/json"
	"fmt"
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/apinodeutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().APINodeRPC().CountAllEnabledAPINodes(this.AdminContext(), &pb.CountAllEnabledAPINodesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	var page = this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	var nodeMaps = []maps.Map{}
	if count > 0 {
		nodesResp, err := this.RPC().APINodeRPC().ListEnabledAPINodes(this.AdminContext(), &pb.ListEnabledAPINodesRequest{
			Offset: page.Offset,
			Size:   page.Size,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		for _, node := range nodesResp.ApiNodes {
			// 状态
			var status = &nodeconfigs.NodeStatus{}
			if len(node.StatusJSON) > 0 {
				err = json.Unmarshal(node.StatusJSON, &status)
				if err != nil {
					logs.Error(err)
					continue
				}
				status.IsActive = status.IsActive && time.Now().Unix()-status.UpdatedAt <= 60 // N秒之内认为活跃
			}

			// Rest地址
			var restAccessAddrs = []string{}
			if node.RestIsOn {
				if len(node.RestHTTPJSON) > 0 {
					httpConfig := &serverconfigs.HTTPProtocolConfig{}
					err = json.Unmarshal(node.RestHTTPJSON, httpConfig)
					if err != nil {
						this.ErrorPage(err)
						return
					}
					_ = httpConfig.Init()
					if httpConfig.IsOn && len(httpConfig.Listen) > 0 {
						for _, listen := range httpConfig.Listen {
							restAccessAddrs = append(restAccessAddrs, listen.FullAddresses()...)
						}
					}
				}

				if len(node.RestHTTPSJSON) > 0 {
					httpsConfig := &serverconfigs.HTTPSProtocolConfig{}
					err = json.Unmarshal(node.RestHTTPSJSON, httpsConfig)
					if err != nil {
						this.ErrorPage(err)
						return
					}
					_ = httpsConfig.Init(context.TODO())
					if httpsConfig.IsOn && len(httpsConfig.Listen) > 0 {
						restAccessAddrs = append(restAccessAddrs, httpsConfig.FullAddresses()...)
					}
				}
			}

			var shouldUpgrade = status.IsActive && len(status.BuildVersion) > 0 && stringutil.VersionCompare(teaconst.APINodeVersion, status.BuildVersion) > 0
			canUpgrade, _ := apinodeutils.CanUpgrade(status.BuildVersion, status.OS, status.Arch)

			nodeMaps = append(nodeMaps, maps.Map{
				"id":              node.Id,
				"isOn":            node.IsOn,
				"name":            node.Name,
				"accessAddrs":     node.AccessAddrs,
				"restAccessAddrs": restAccessAddrs,
				"isPrimary":       node.IsPrimary,
				"status": maps.Map{
					"isActive":      status.IsActive,
					"updatedAt":     status.UpdatedAt,
					"hostname":      status.Hostname,
					"cpuUsage":      status.CPUUsage,
					"cpuUsageText":  fmt.Sprintf("%.2f%%", status.CPUUsage*100),
					"memUsage":      status.MemoryUsage,
					"memUsageText":  fmt.Sprintf("%.2f%%", status.MemoryUsage*100),
					"buildVersion":  status.BuildVersion,
					"latestVersion": teaconst.APINodeVersion,
					"shouldUpgrade": shouldUpgrade,
					"canUpgrade":    shouldUpgrade && canUpgrade,
				},
			})
		}
	}
	this.Data["nodes"] = nodeMaps

	// 检查是否有调试数据
	countMethodStatsResp, err := this.RPC().APIMethodStatRPC().CountAPIMethodStatsWithDay(this.AdminContext(), &pb.CountAPIMethodStatsWithDayRequest{Day: timeutil.Format("Ymd")})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasMethodStats"] = countMethodStatsResp.Count > 0

	this.Show()
}
