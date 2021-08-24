// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package tasks

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "task")
}

func (this *IndexAction) RunGet(params struct {
	Status int32
}) {
	this.Data["status"] = params.Status
	if params.Status > 3 {
		params.Status = 0
	}

	countWaitingResp, err := this.RPC().MessageTaskRPC().CountMessageTasksWithStatus(this.AdminContext(), &pb.CountMessageTasksWithStatusRequest{Status: pb.CountMessageTasksWithStatusRequest_MessageTaskStatusNone})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var countWaiting = countWaitingResp.Count
	this.Data["countWaiting"] = countWaiting

	countFailedResp, err := this.RPC().MessageTaskRPC().CountMessageTasksWithStatus(this.AdminContext(), &pb.CountMessageTasksWithStatusRequest{Status: pb.CountMessageTasksWithStatusRequest_MessageTaskStatusFailed})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var countFailed = countFailedResp.Count
	this.Data["countFailed"] = countFailed

	// 列表
	var total = int64(0)
	switch params.Status {
	case 0:
		total = countWaiting
	case 3:
		total = countFailed
	}
	page := this.NewPage(total)
	this.Data["page"] = page.AsHTML()

	var taskMaps = []maps.Map{}
	tasksResp, err := this.RPC().MessageTaskRPC().ListMessageTasksWithStatus(this.AdminContext(), &pb.ListMessageTasksWithStatusRequest{
		Status: pb.ListMessageTasksWithStatusRequest_Status(params.Status),
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	for _, task := range tasksResp.MessageTasks {
		var resultMap = maps.Map{}
		var result = task.Result
		if result != nil {
			resultMap = maps.Map{
				"isOk":     result.IsOk,
				"error":    result.Error,
				"response": result.Response,
			}
		}

		//var recipients = []string{}
		var user = ""
		var instanceMap maps.Map
		if task.MessageRecipient != nil {
			user = task.MessageRecipient.User
			if task.MessageRecipient.MessageMediaInstance != nil {
				instanceMap = maps.Map{
					"id":   task.MessageRecipient.MessageMediaInstance.Id,
					"name": task.MessageRecipient.MessageMediaInstance.Name,
				}
			}
		}

		taskMaps = append(taskMaps, maps.Map{
			"id":          task.Id,
			"subject":     task.Subject,
			"body":        task.Body,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", task.CreatedAt),
			"result":      resultMap,
			"status":      task.Status,
			"user":        user,
			"instance":    instanceMap,
		})
	}
	this.Data["tasks"] = taskMaps

	this.Show()
}
