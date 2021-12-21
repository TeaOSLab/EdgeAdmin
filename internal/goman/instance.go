// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package goman

import "time"

type Instance struct {
	Id          uint64
	CreatedTime time.Time
	File        string
	Line        int
}
