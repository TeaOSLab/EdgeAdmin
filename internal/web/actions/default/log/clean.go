package log

import (
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configloaders"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type CleanAction struct {
	actionutils.ParentAction
}

func (this *CleanAction) Init() {
	this.Nav("", "", "clean")
}

func (this *CleanAction) RunGet(params struct{}) {
	// 读取配置
	config, err := configloaders.LoadLogConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if !config.CanClean {
		this.WriteString("已设置不能清理")
		return
	}
	this.Data["logConfig"] = config

	sizeResp, err := this.RPC().LogRPC().SumLogsSize(this.AdminContext(), &pb.SumLogsSizeRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	sizeHuman := ""
	if sizeResp.SizeBytes < 1024 {
		sizeHuman = numberutils.FormatInt64(sizeResp.SizeBytes) + "字节"
	} else if sizeResp.SizeBytes < 1024*1024 {
		sizeHuman = fmt.Sprintf("%.2fK", float64(sizeResp.SizeBytes)/1024)
	} else if sizeResp.SizeBytes < 1024*1024*1024 {
		sizeHuman = fmt.Sprintf("%.2fM", float64(sizeResp.SizeBytes)/1024/1024)
	} else {
		sizeHuman = fmt.Sprintf("%.2fG", float64(sizeResp.SizeBytes)/1024/1024/1024)
	}
	this.Data["size"] = sizeHuman

	this.Show()
}

func (this *CleanAction) RunPost(params struct {
	Type string
	Days int32

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// 读取配置
	config, err := configloaders.LoadLogConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if !config.CanClean {
		this.WriteString("已设置不能清理")
		return
	}

	switch params.Type {
	case "all":
		defer this.CreateLogInfo(codes.Log_LogCleanAllLogs)

		_, err := this.RPC().LogRPC().CleanLogsPermanently(this.AdminContext(), &pb.CleanLogsPermanentlyRequest{
			Days:     0,
			ClearAll: true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	case "days":
		defer this.CreateLogInfo(codes.Log_LogCleanLogsDaysBefore, params.Days)

		_, err := this.RPC().LogRPC().CleanLogsPermanently(this.AdminContext(), &pb.CleanLogsPermanentlyRequest{
			Days:     params.Days,
			ClearAll: false,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
	default:
		this.Fail("不支持的清理方式 '" + params.Type + "'")
	}

	this.Success()
}
