// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"sort"
)

type TaskAction struct {
	actionutils.ParentAction
}

func (this *TaskAction) Init() {
	this.Nav("", "", "task")
}

func (this *TaskAction) RunGet(params struct {
	TaskId int64
}) {
	// 初始化菜单数据
	err := InitMenu(this.Parent())
	if err != nil {
		this.ErrorPage(err)
		return
	}

	taskResp, err := this.RPC().HTTPCacheTaskRPC().FindEnabledHTTPCacheTask(this.AdminContext(), &pb.FindEnabledHTTPCacheTaskRequest{HttpCacheTaskId: params.TaskId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var task = taskResp.HttpCacheTask
	if task == nil {
		this.NotFound("HTTPCacheTask", params.TaskId)
		return
	}

	// 用户
	var userMap = maps.Map{"id": 0, "username": "", "fullname": ""}
	if task.User != nil {
		userMap = maps.Map{
			"id":       task.User.Id,
			"username": task.User.Username,
			"fullname": task.User.Fullname,
		}
	}

	// keys
	var keyMaps = []maps.Map{}
	for _, key := range task.HttpCacheTaskKeys {
		// 错误信息
		var errorMaps = []maps.Map{}

		if len(key.ErrorsJSON) > 0 {
			var m = map[int64]string{}
			err = json.Unmarshal(key.ErrorsJSON, &m)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			for nodeId, errString := range m {
				errorMaps = append(errorMaps, maps.Map{
					"nodeId": nodeId,
					"error":  errString,
				})
			}
		}

		// 错误信息排序
		if len(errorMaps) > 0 {
			sort.Slice(errorMaps, func(i, j int) bool {
				var m1 = errorMaps[i]
				var m2 = errorMaps[j]

				return m1.GetInt64("nodeId") < m2.GetInt64("nodeId")
			})
		}

		keyMaps = append(keyMaps, maps.Map{
			"key":    key.Key,
			"isDone": key.IsDone,
			"errors": errorMaps,
		})
	}

	this.Data["task"] = maps.Map{
		"id":          task.Id,
		"type":        task.Type,
		"keyType":     task.KeyType,
		"createdTime": timeutil.FormatTime("Y-m-d H:i:s", task.CreatedAt),
		"doneTime":    timeutil.FormatTime("Y-m-d H:i:s", task.DoneAt),
		"isDone":      task.IsDone,
		"isOk":        task.IsOk,
		"keys":        keyMaps,
		"user":        userMap,
	}

	this.Show()
}
