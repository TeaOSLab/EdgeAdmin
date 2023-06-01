// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !plus

package https

func (this *IndexAction) checkSupportsHTTP3(clusterId int64) (bool, error) {
	return false, nil
}
