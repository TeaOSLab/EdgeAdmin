package group

import (	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	GroupId int64
}) {
	// 创建日志
	defer this.CreateLogInfo(codes.ServerGroup_LogDeleteServerGroup, params.GroupId)

	// 检查是否正在使用
	countResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithServerGroupId(this.AdminContext(), &pb.CountAllEnabledServersWithServerGroupIdRequest{ServerGroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if countResp.Count > 0 {
		this.Fail("此分组正在被使用不能删除，请修改相关服务后再删除")
	}

	_, err = this.RPC().ServerGroupRPC().DeleteServerGroup(this.AdminContext(), &pb.DeleteServerGroupRequest{ServerGroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
