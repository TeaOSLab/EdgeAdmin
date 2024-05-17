// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.
//go:build !plus

package node

func (this *UpdateAction) CanUpdateLevel(level int32) bool {
	return level <= 1
}
