// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package files

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/types"
	"mime"
	"path/filepath"
)

type FileAction struct {
	actionutils.ParentAction
}

func (this *FileAction) Init() {
	this.Nav("", "", "")
}

func (this *FileAction) RunGet(params struct {
	FileId int64
}) {
	fileResp, err := this.RPC().FileRPC().FindEnabledFile(this.AdminContext(), &pb.FindEnabledFileRequest{FileId: params.FileId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var file = fileResp.File
	if file == nil {
		this.NotFound("File", params.FileId)
		return
	}

	chunkIdsResp, err := this.RPC().FileChunkRPC().FindAllFileChunkIds(this.AdminContext(), &pb.FindAllFileChunkIdsRequest{FileId: file.Id})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.AddHeader("Content-Length", types.String(file.Size))
	if len(file.MimeType) > 0 {
		this.AddHeader("Content-Type", file.MimeType)
	} else if len(file.Filename) > 0 {
		var ext = filepath.Ext(file.Filename)
		var mimeType = mime.TypeByExtension(ext)
		this.AddHeader("Content-Type", mimeType)
	}

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
