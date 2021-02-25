package userui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"io"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	config, err := configloaders.LoadUserUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["config"] = config

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ProductName        string
	UserSystemName     string
	ShowOpenSourceInfo bool
	ShowVersion        bool
	Version            string
	ShowFinance        bool
	FaviconFile        *actions.File
	LogoFile           *actions.File

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	params.Must.
		Field("productName", params.ProductName).
		Require("请输入产品名称").
		Field("userSystemName", params.UserSystemName).
		Require("请输入管理员系统名称")

	config, err := configloaders.LoadUserUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	config.ProductName = params.ProductName
	config.UserSystemName = params.UserSystemName
	config.ShowOpenSourceInfo = params.ShowOpenSourceInfo
	config.ShowVersion = params.ShowVersion
	config.Version = params.Version
	config.ShowFinance = params.ShowFinance

	// 上传Favicon文件
	if params.FaviconFile != nil {
		createResp, err := this.RPC().FileRPC().CreateFile(this.AdminContext(), &pb.CreateFileRequest{
			Filename: params.FaviconFile.Filename,
			Size:     params.FaviconFile.Size,
			IsPublic: true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		fileId := createResp.FileId

		// 上传内容
		buf := make([]byte, 512*1024)
		reader, err := params.FaviconFile.OriginFile.Open()
		if err != nil {
			this.ErrorPage(err)
			return
		}
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
		}
		config.FaviconFileId = fileId
	}

	// 上传Logo文件
	if params.LogoFile != nil {
		createResp, err := this.RPC().FileRPC().CreateFile(this.AdminContext(), &pb.CreateFileRequest{
			Filename: params.LogoFile.Filename,
			Size:     params.LogoFile.Size,
			IsPublic: true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		fileId := createResp.FileId

		// 上传内容
		buf := make([]byte, 512*1024)
		reader, err := params.LogoFile.OriginFile.Open()
		if err != nil {
			this.ErrorPage(err)
			return
		}
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
		}
		config.LogoFileId = fileId
	}

	err = configloaders.UpdateUserUIConfig(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
