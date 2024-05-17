// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package iplists

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/tealeg/xlsx/v3"
	"strconv"
)

type ExportDataAction struct {
	actionutils.ParentAction
}

func (this *ExportDataAction) Init() {
	this.Nav("", "", "")
}

func (this *ExportDataAction) RunGet(params struct {
	ListId int64
	Format string
}) {
	defer this.CreateLogInfo(codes.IPList_LogExportIPList, params.ListId)

	var ext string
	var jsonMaps = []maps.Map{}
	var xlsxFile *xlsx.File
	var xlsxSheet *xlsx.Sheet
	var csvWriter *csv.Writer
	var csvBuffer *bytes.Buffer

	var data []byte

	switch params.Format {
	case "xlsx":
		ext = ".xlsx"
		xlsxFile = xlsx.NewFile()
		var err error
		xlsxSheet, err = xlsxFile.AddSheet("IP名单")
		if err != nil {
			this.ErrorPage(err)
			return
		}

		var row = xlsxSheet.AddRow()
		row.SetHeight(26)
		row.AddCell().SetValue("IP/IP段")
		row.AddCell().SetValue("过期时间戳")
		row.AddCell().SetValue("类型")
		row.AddCell().SetValue("级别")
		row.AddCell().SetValue("备注")
	case "csv":
		ext = ".csv"
		csvBuffer = &bytes.Buffer{}
		csvWriter = csv.NewWriter(csvBuffer)
	case "txt":
		ext = ".txt"
	case "json":
		ext = ".json"
	default:
		this.WriteString("请选择正确的导出格式")
		return
	}

	var offset int64 = 0
	var size int64 = 1000
	for {
		itemsResp, err := this.RPC().IPItemRPC().ListIPItemsWithListId(this.AdminContext(), &pb.ListIPItemsWithListIdRequest{
			IpListId: params.ListId,
			Offset:   offset,
			Size:     size,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(itemsResp.IpItems) == 0 {
			break
		}
		for _, item := range itemsResp.IpItems {
			switch params.Format {
			case "xlsx":
				var row = xlsxSheet.AddRow()
				row.SetHeight(26)
				row.AddCell().SetValue(item.Value)
				row.AddCell().SetValue(types.String(item.ExpiredAt))
				row.AddCell().SetValue(item.Type)
				row.AddCell().SetValue(item.EventLevel)
				row.AddCell().SetValue(item.Reason)
			case "csv":
				err = csvWriter.Write([]string{item.Value, types.String(item.ExpiredAt), item.Type, item.EventLevel, item.Reason})
				if err != nil {
					this.ErrorPage(err)
					return
				}
			case "txt":
				data = append(data, item.Value+","+types.String(item.ExpiredAt)+","+item.Type+","+item.EventLevel+","+item.Reason...)
				data = append(data, '\n')
			case "json":
				jsonMaps = append(jsonMaps, maps.Map{
					"value":      item.Value,
					"expiredAt":  item.ExpiredAt,
					"type":       item.Type,
					"eventLevel": item.EventLevel,
					"reason":     item.Reason,
				})
			}
		}
		offset += size
	}

	switch params.Format {
	case "xlsx":
		var buf = &bytes.Buffer{}
		err := xlsxFile.Write(buf)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		data = buf.Bytes()
	case "csv":
		csvWriter.Flush()
		data = csvBuffer.Bytes()
	case "json":
		var err error
		data, err = json.Marshal(jsonMaps)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.AddHeader("Content-Disposition", "attachment; filename=\"ip-list-"+numberutils.FormatInt64(params.ListId)+ext+"\";")
	this.AddHeader("Content-Length", strconv.Itoa(len(data)))
	_, _ = this.Write(data)
}
