// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "node", "update")
	this.SecondMenu("cache")
}

func (this *IndexAction) RunGet(params struct {
	NodeId int64
}) {
	node, err := nodeutils.InitNodeInfo(this.Parent(), params.NodeId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 缓存硬盘 & 内存容量
	var maxCacheDiskCapacity maps.Map
	if node.MaxCacheDiskCapacity != nil {
		maxCacheDiskCapacity = maps.Map{
			"count": node.MaxCacheDiskCapacity.Count,
			"unit":  node.MaxCacheDiskCapacity.Unit,
		}
	} else {
		maxCacheDiskCapacity = maps.Map{
			"count": 0,
			"unit":  "gb",
		}
	}

	var maxCacheMemoryCapacity maps.Map
	if node.MaxCacheMemoryCapacity != nil {
		maxCacheMemoryCapacity = maps.Map{
			"count": node.MaxCacheMemoryCapacity.Count,
			"unit":  node.MaxCacheMemoryCapacity.Unit,
		}
	} else {
		maxCacheMemoryCapacity = maps.Map{
			"count": 0,
			"unit":  "gb",
		}
	}

	var diskSubDirs = []*serverconfigs.CacheDir{}
	if len(node.CacheDiskSubDirsJSON) > 0 {
		err = json.Unmarshal(node.CacheDiskSubDirsJSON, &diskSubDirs)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	var nodeMap = this.Data["node"].(maps.Map)
	nodeMap["maxCacheDiskCapacity"] = maxCacheDiskCapacity
	nodeMap["cacheDiskDir"] = node.CacheDiskDir
	nodeMap["cacheDiskSubDirs"] = diskSubDirs
	nodeMap["maxCacheMemoryCapacity"] = maxCacheMemoryCapacity

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	NodeId                     int64
	MaxCacheDiskCapacityJSON   []byte
	CacheDiskDir               string
	CacheDiskSubDirsJSON       []byte
	MaxCacheMemoryCapacityJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.NodeCache_LogUpdateNodeCacheSettings, params.NodeId)

	// 缓存硬盘 & 内存容量
	var pbMaxCacheDiskCapacity *pb.SizeCapacity
	if len(params.MaxCacheDiskCapacityJSON) > 0 {
		var sizeCapacity = &shared.SizeCapacity{}
		err := json.Unmarshal(params.MaxCacheDiskCapacityJSON, sizeCapacity)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		pbMaxCacheDiskCapacity = &pb.SizeCapacity{
			Count: sizeCapacity.Count,
			Unit:  sizeCapacity.Unit,
		}
	}

	var pbMaxCacheMemoryCapacity *pb.SizeCapacity
	if len(params.MaxCacheMemoryCapacityJSON) > 0 {
		var sizeCapacity = &shared.SizeCapacity{}
		err := json.Unmarshal(params.MaxCacheMemoryCapacityJSON, sizeCapacity)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		pbMaxCacheMemoryCapacity = &pb.SizeCapacity{
			Count: sizeCapacity.Count,
			Unit:  sizeCapacity.Unit,
		}
	}

	if len(params.CacheDiskSubDirsJSON) > 0 {
		var cacheSubDirs = []*serverconfigs.CacheDir{}
		err := json.Unmarshal(params.CacheDiskSubDirsJSON, &cacheSubDirs)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	_, err := this.RPC().NodeRPC().UpdateNodeCache(this.AdminContext(), &pb.UpdateNodeCacheRequest{
		NodeId:                 params.NodeId,
		MaxCacheDiskCapacity:   pbMaxCacheDiskCapacity,
		CacheDiskDir:           params.CacheDiskDir,
		CacheDiskSubDirsJSON:   params.CacheDiskSubDirsJSON,
		MaxCacheMemoryCapacity: pbMaxCacheMemoryCapacity,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
