package iplibrary

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DownloadAction struct {
	actionutils.ParentAction
}

func (this *DownloadAction) Init() {
	this.Nav("", "", "")
}

func (this *DownloadAction) RunGet(params struct {
	LibraryId int64
}) {
	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "下载IP库 %d", params.LibraryId)

	libraryResp, err := this.RPC().IPLibraryRPC().FindEnabledIPLibrary(this.AdminContext(), &pb.FindEnabledIPLibraryRequest{IpLibraryId: params.LibraryId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if libraryResp.IpLibrary == nil || libraryResp.IpLibrary.File == nil {
		this.NotFound("ipLibrary", params.LibraryId)
		return
	}

	file := libraryResp.IpLibrary.File
	chunkIdsResp, err := this.RPC().FileChunkRPC().FindAllFileChunkIds(this.AdminContext(), &pb.FindAllFileChunkIdsRequest{FileId: file.Id})
	if err != nil {
		this.ErrorPage(err)
	}

	this.AddHeader("Content-Disposition", "attachment; filename=\""+file.Filename+"\";")
	for _, chunkId := range chunkIdsResp.FileChunkIds {
		chunkResp, err := this.RPC().FileChunkRPC().DownloadFileChunk(this.AdminContext(), &pb.DownloadFileChunkRequest{FileChunkId: chunkId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if chunkResp.FileChunk != nil {
			this.Write(chunkResp.FileChunk.Data)
		}
	}
}
