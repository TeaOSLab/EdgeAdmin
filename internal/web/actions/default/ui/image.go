package ui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"mime"
	"path/filepath"
	"strconv"
)

// 公开的图片，不需要检查用户权限
type ImageAction struct {
	actionutils.ParentAction
}

func (this *ImageAction) Init() {
	this.Nav("", "", "")
}

func (this *ImageAction) RunGet(params struct {
	FileId int64
}) {
	fileResp, err := this.RPC().FileRPC().FindEnabledFile(this.AdminContext(), &pb.FindEnabledFileRequest{FileId: params.FileId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	file := fileResp.File
	if file == nil {
		this.NotFound("file", params.FileId)
		return
	}

	if !file.IsPublic {
		this.NotFound("file", params.FileId)
		return
	}

	chunkIdsResp, err := this.RPC().FileChunkRPC().FindAllFileChunkIds(this.AdminContext(), &pb.FindAllFileChunkIdsRequest{FileId: file.Id})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	mimeType := ""
	if len(file.Filename) > 0 {
		ext := filepath.Ext(file.Filename)
		mimeType = mime.TypeByExtension(ext)
	}
	if len(mimeType) == 0 {
		mimeType = "image/png"
	}

	this.AddHeader("Last-Modified", "Fri, 06 Sep 2019 08:29:50 GMT")
	this.AddHeader("Content-Type", mimeType)
	this.AddHeader("Content-Length", strconv.FormatInt(file.Size, 10))
	for _, chunkId := range chunkIdsResp.FileChunkIds {
		chunkResp, err := this.RPC().FileChunkRPC().DownloadFileChunk(this.AdminContext(), &pb.DownloadFileChunkRequest{FileChunkId: chunkId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if chunkResp.FileChunk == nil {
			continue
		}
		_, _ = this.Write(chunkResp.FileChunk.Data)
	}
}
