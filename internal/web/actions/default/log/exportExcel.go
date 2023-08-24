package log

import (
	"bytes"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/iplibrary"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"github.com/tealeg/xlsx/v3"
	"strconv"
)

type ExportExcelAction struct {
	actionutils.ParentAction
}

func (this *ExportExcelAction) Init() {
	this.Nav("", "", "")
}

func (this *ExportExcelAction) RunGet(params struct {
	DayFrom  string
	DayTo    string
	Keyword  string
	UserType string
	Level    string
}) {
	logsResp, err := this.RPC().LogRPC().ListLogs(this.AdminContext(), &pb.ListLogsRequest{
		Offset:   0,
		Size:     10000, // 日志最大导出10000条，TODO 将来可以配置
		DayFrom:  params.DayFrom,
		DayTo:    params.DayTo,
		Keyword:  params.Keyword,
		UserType: params.UserType,
		Level:    params.Level,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var wb = xlsx.NewFile()
	sheet, err := wb.AddSheet("default")
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 头部
	{
		var row = sheet.AddRow()
		row.SetHeight(25)
		row.AddCell().SetString("ID")
		row.AddCell().SetString("日期")
		row.AddCell().SetString("用户")
		row.AddCell().SetString("描述")
		row.AddCell().SetString("IP")
		row.AddCell().SetString("区域")
		row.AddCell().SetString("运营商")
		row.AddCell().SetString("页面地址")
		row.AddCell().SetString("级别")
	}

	// 数据
	for _, log := range logsResp.Logs {
		var regionName = ""
		var ispName = ""

		var ipRegion = iplibrary.LookupIP(log.Ip)
		if ipRegion != nil && ipRegion.IsOk() {
			regionName = ipRegion.RegionSummary()
			ispName = ipRegion.ProviderName()
		}

		var row = sheet.AddRow()
		row.SetHeight(25)
		row.AddCell().SetInt64(log.Id)
		row.AddCell().SetString(timeutil.FormatTime("Y-m-d H:i:s", log.CreatedAt))
		if log.UserId > 0 {
			row.AddCell().SetString("用户 | " + log.UserName)
		} else {
			row.AddCell().SetString(log.UserName)
		}
		row.AddCell().SetString(log.Description)
		row.AddCell().SetString(log.Ip)
		row.AddCell().SetString(regionName)
		row.AddCell().SetString(ispName)
		row.AddCell().SetString(log.Action)

		var levelName = ""
		switch log.Level {
		case "info":
			levelName = "信息"
		case "warn", "warning":
			levelName = "警告"
		case "error":
			levelName = "错误"
		}
		row.AddCell().SetString(levelName)
	}

	this.AddHeader("Content-Type", "application/vnd.ms-excel")
	this.AddHeader("Content-Disposition", "attachment; filename=\"LOG-"+timeutil.Format("YmdHis")+".xlsx\"")
	this.AddHeader("Cache-Control", "max-age=0")

	var buf = bytes.NewBuffer([]byte{})
	err = wb.Write(buf)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.AddHeader("Content-Length", strconv.Itoa(buf.Len()))
	_, _ = this.Write(buf.Bytes())
}
