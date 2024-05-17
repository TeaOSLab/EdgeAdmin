// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !plus

package dashboard

func (this *IndexAction) checkPlus() bool {
	return false
}
