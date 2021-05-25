// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package clusters

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().NSClusterRPC().CountAllEnabledNSClusters(this.AdminContext(), &pb.CountAllEnabledNSClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	clustersResp, err := this.RPC().NSClusterRPC().ListEnabledNSClusters(this.AdminContext(), &pb.ListEnabledNSClustersRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	clusterMaps := []maps.Map{}
	for _, cluster := range clustersResp.NsClusters {
		clusterMaps = append(clusterMaps, maps.Map{
			"id":   cluster.Id,
			"name": cluster.Name,
			"isOn": cluster.IsOn,
		})
	}
	this.Data["clusters"] = clusterMaps

	this.Show()
}
