package iplibrary

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"io"
)

type UploadPopupAction struct {
	actionutils.ParentAction
}

func (this *UploadPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UploadPopupAction) RunGet(params struct{}) {
	this.Data["types"] = serverconfigs.IPLibraryTypes

	this.Show()
}

func (this *UploadPopupAction) RunPost(params struct {
	Type string
	File *actions.File

	Must *actions.Must
}) {
	libraryType := serverconfigs.FindIPLibraryWithType(params.Type)
	if libraryType == nil {
		this.Fail("错误的IP类型")
	}

	if params.File == nil {
		this.Fail("请选择要上传的文件")
	}

	if params.File.Size == 0 {
		this.Fail("文件内容不能为空")
	}

	if params.File.Ext != libraryType.GetString("ext") {
		this.Fail("IP库文件扩展名错误，应该为：" + libraryType.GetString("ext"))
	}

	reader, err := params.File.OriginFile.Open()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	defer func() {
		_ = reader.Close()
	}()

	// 创建文件
	fileResp, err := this.RPC().FileRPC().CreateFile(this.AdminContext(), &pb.CreateFileRequest{
		Filename: params.File.Filename,
		Size:     params.File.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	fileId := fileResp.FileId

	// 上传内容
	buf := make([]byte, 512*1024)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			_, err = this.RPC().FileChunkRPC().CreateFileChunk(this.AdminContext(), &pb.CreateFileChunkRequest{
				FileId: fileId,
				Data:   buf[:n],
			})
			if err != nil {
				this.Fail("上传失败：" + err.Error())
			}
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			this.Fail("上传失败：" + err.Error())
		}
	}

	// 置为已完成
	_, err = this.RPC().FileRPC().UpdateFileFinished(this.AdminContext(), &pb.UpdateFileFinishedRequest{FileId: fileId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 保存
	createResp, err := this.RPC().IPLibraryRPC().CreateIPLibrary(this.AdminContext(), &pb.CreateIPLibraryRequest{
		Type:   params.Type,
		FileId: fileId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "上传IP库 %d", createResp.IpLibraryId)

	this.Success()
}
