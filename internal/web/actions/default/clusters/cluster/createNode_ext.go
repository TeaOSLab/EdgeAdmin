// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !plus

package cluster

func (this *CreateNodeAction) findNodesQuota() (maxNodes int32, leftNodes int32, err error) {
	return
}
