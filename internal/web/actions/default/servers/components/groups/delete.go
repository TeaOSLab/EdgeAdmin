package groups

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	GroupId int64
}) {
	// 检查是否正在使用
	countResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithGroupId(this.AdminContext(), &pb.CountAllEnabledServersWithGroupIdRequest{GroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if countResp.Count > 0 {
		this.Fail("此分组正在被使用不能删除，请修改相关服务后再删除")
	}

	_, err = this.RPC().ServerGroupRPC().DeleteServerGroup(this.AdminContext(), &pb.DeleteServerGroupRequest{GroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
