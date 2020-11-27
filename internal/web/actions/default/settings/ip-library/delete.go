package iplibrary

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	LibraryId int64
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "删除IP库 %d", params.LibraryId)

	_, err := this.RPC().IPLibraryRPC().DeleteIPLibrary(this.AdminContext(), &pb.DeleteIPLibraryRequest{IpLibraryId: params.LibraryId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
