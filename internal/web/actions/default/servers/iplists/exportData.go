// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package iplists

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
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
	defer this.CreateLogInfo("导出IP名单 %d", params.ListId)

	var err error
	var ext = ""
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
		if err != nil {
			this.ErrorPage(err)
			return
		}
		xlsxSheet, err = xlsxFile.AddSheet("IP名单")
		if err != nil {
			this.ErrorPage(err)
			return
		}

		row := xlsxSheet.AddRow()
		row.SetHeight(26)
		row.AddCell().SetValue("开始IP")
		row.AddCell().SetValue("结束IP")
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
				row := xlsxSheet.AddRow()
				row.SetHeight(26)
				row.AddCell().SetValue(item.IpFrom)
				row.AddCell().SetValue(item.IpTo)
				row.AddCell().SetValue(types.String(item.ExpiredAt))
				row.AddCell().SetValue(item.Type)
				row.AddCell().SetValue(item.EventLevel)
				row.AddCell().SetValue(item.Reason)
			case "csv":
				err = csvWriter.Write([]string{item.IpFrom, item.IpTo, types.String(item.ExpiredAt), item.Type, item.EventLevel, item.Reason})
				if err != nil {
					this.ErrorPage(err)
					return
				}
			case "txt":
				data = append(data, item.IpFrom+","+item.IpTo+","+types.String(item.ExpiredAt)+","+item.Type+","+item.EventLevel+","+item.Reason...)
				data = append(data, '\n')
			case "json":
				jsonMaps = append(jsonMaps, maps.Map{
					"ipFrom":     item.IpFrom,
					"ipTo":       item.IpTo,
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
		err = xlsxFile.Write(buf)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		data = buf.Bytes()
	case "csv":
		csvWriter.Flush()
		data = csvBuffer.Bytes()
	case "json":
		data, err = json.Marshal(jsonMaps)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.AddHeader("Content-Disposition", "attachment; filename=\"ip-list-"+numberutils.FormatInt64(params.ListId)+ext+"\";")
	this.AddHeader("Content-Length", strconv.Itoa(len(data)))
	this.Write(data)
}
