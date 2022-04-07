// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved.
//go:build !plus
// +build !plus

package node

func (this *UpdateAction) CanUpdateLevel(level int32) bool {
	return level <= 1
}
