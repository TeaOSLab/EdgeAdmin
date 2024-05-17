// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type TasksAction struct {
	actionutils.ParentAction
}

func (this *TasksAction) Init() {
	this.Nav("", "", "task")
}

func (this *TasksAction) RunGet(params struct{}) {
	// 初始化菜单数据
	err := InitMenu(this.Parent())
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 任务数量
	countResp, err := this.RPC().HTTPCacheTaskRPC().CountHTTPCacheTasks(this.AdminContext(), &pb.CountHTTPCacheTasksRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	var page = this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	// 任务列表
	var taskMaps = []maps.Map{}
	tasksResp, err := this.RPC().HTTPCacheTaskRPC().ListHTTPCacheTasks(this.AdminContext(), &pb.ListHTTPCacheTasksRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, task := range tasksResp.HttpCacheTasks {
		var userMap = maps.Map{"id": 0, "username": "", "fullname": ""}
		if task.User != nil {
			userMap = maps.Map{
				"id":       task.User.Id,
				"username": task.User.Username,
				"fullname": task.User.Fullname,
			}
		}

		taskMaps = append(taskMaps, maps.Map{
			"id":          task.Id,
			"type":        task.Type,
			"keyType":     task.KeyType,
			"isDone":      task.IsDone,
			"isOk":        task.IsOk,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", task.CreatedAt),
			"description": task.Description,
			"user":        userMap,
		})
	}

	this.Data["tasks"] = taskMaps

	this.Show()
}
