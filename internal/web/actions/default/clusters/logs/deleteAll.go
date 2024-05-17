// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package logs

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAllAction struct {
	actionutils.ParentAction
}

func (this *DeleteAllAction) RunPost(params struct {
	DayFrom   string
	DayTo     string
	Keyword   string
	Level     string
	Type      string // unread, needFix
	Tag       string
	ClusterId int64
	NodeId    int64
}) {
	defer this.CreateLogInfo(codes.NodeLog_LogDeleteNodeLogsBatch)

	// 目前仅允许通过关键词删除，防止误删
	if len(params.Keyword) == 0 {
		this.Fail("目前仅允许通过关键词删除")
		return
	}

	var fixedState configutils.BoolState = 0
	var allServers = false
	if params.Type == "needFix" {
		fixedState = configutils.BoolStateNo
		allServers = true
	}

	_, err := this.RPC().NodeLogRPC().DeleteNodeLogs(this.AdminContext(), &pb.DeleteNodeLogsRequest{
		NodeClusterId: params.ClusterId,
		NodeId:        params.NodeId,
		Role:          nodeconfigs.NodeRoleNode,
		DayFrom:       params.DayFrom,
		DayTo:         params.DayTo,
		Keyword:       params.Keyword,
		Level:         params.Level,
		IsUnread:      params.Type == "unread",
		Tag:           params.Tag,
		FixedState:    int32(fixedState),
		AllServers:    allServers,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
