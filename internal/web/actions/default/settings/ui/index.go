package ui

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/userconfigs"
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
	config, err := configloaders.LoadAdminUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if config.DefaultPageSize == 0 {
		config.DefaultPageSize = 10
	}
	this.Data["config"] = config

	// 时区
	this.Data["timeZoneGroups"] = nodeconfigs.FindAllTimeZoneGroups()
	this.Data["timeZoneLocations"] = nodeconfigs.FindAllTimeZoneLocations()

	if len(config.TimeZone) == 0 {
		config.TimeZone = nodeconfigs.DefaultTimeZoneLocation
	}
	this.Data["timeZoneLocation"] = nodeconfigs.FindTimeZoneLocation(config.TimeZone)

	this.filterConfig(config)

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ProductName        string
	AdminSystemName    string
	ShowOpenSourceInfo bool
	ShowFinance        bool
	ShowVersion        bool
	Version            string
	FaviconFile        *actions.File
	LogoFile           *actions.File
	DefaultPageSize    int
	TimeZone           string
	DnsResolverType    string

	SupportModuleCDN bool
	SupportModuleNS  bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.AdminUI_LogUpdateUISettings)

	params.Must.
		Field("productName", params.ProductName).
		Require("请输入产品名称").
		Field("adminSystemName", params.AdminSystemName).
		Require("请输入管理员系统名称").
		Field("defaultPageSize", params.DefaultPageSize).
		Gte(0, "默认每页显示数不能小于0").
		Lte(100, "默认每页显示数不能大于100")

	config, err := configloaders.LoadAdminUIConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	config.ProductName = params.ProductName
	config.AdminSystemName = params.AdminSystemName
	config.ShowOpenSourceInfo = params.ShowOpenSourceInfo
	config.ShowFinance = params.ShowFinance
	config.ShowVersion = params.ShowVersion
	config.Version = params.Version
	config.TimeZone = params.TimeZone
	config.DNSResolver.Type = params.DnsResolverType

	if params.DefaultPageSize > 0 {
		config.DefaultPageSize = params.DefaultPageSize
	} else {
		config.DefaultPageSize = 10
	}

	config.Modules = []userconfigs.UserModule{}
	if params.SupportModuleCDN {
		config.Modules = append(config.Modules, userconfigs.UserModuleCDN)
	}
	if params.SupportModuleNS {
		config.Modules = append(config.Modules, userconfigs.UserModuleNS)
	}

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
		var fileId = createResp.FileId

		// 上传内容
		var buf = make([]byte, 512*1024)
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
		var fileId = createResp.FileId

		// 上传内容
		var buf = make([]byte, 512*1024)
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

	err = configloaders.UpdateAdminUIConfig(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
