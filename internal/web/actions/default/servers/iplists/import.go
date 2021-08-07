// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package iplists

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/types"
	"github.com/tealeg/xlsx/v3"
	"io"
	"net"
	"regexp"
	"strings"
)

type ImportAction struct {
	actionutils.ParentAction
}

func (this *ImportAction) Init() {
	this.Nav("", "", "import")
}

func (this *ImportAction) RunGet(params struct {
	ListId int64
}) {
	err := InitIPList(this.Parent(), params.ListId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Show()
}

func (this *ImportAction) RunPost(params struct {
	ListId int64
	File   *actions.File

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("导入IP名单 %d", params.ListId)

	existsResp, err := this.RPC().IPListRPC().ExistsEnabledIPList(this.AdminContext(), &pb.ExistsEnabledIPListRequest{IpListId: params.ListId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if !existsResp.Exists {
		this.Fail("IP名单不存在")
	}

	if params.File == nil {
		this.Fail("请选择要导入的IP文件")
	}

	// 检查文件扩展名
	if !regexp.MustCompile(`(?i)\.(xlsx|csv|json|txt)$`).MatchString(params.File.Filename) {
		this.Fail("不支持当前格式的文件导入")
	}

	var ext = strings.ToLower(params.File.Ext)

	data, err := params.File.Read()
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var countIgnore = 0
	var items = []*pb.IPItem{}
	switch ext {
	case ".xlsx":
		file, err := xlsx.OpenBinary(data)
		if err != nil {
			this.Fail("Excel读取错误：" + err.Error())
		}
		if len(file.Sheets) > 0 {
			var sheet = file.Sheets[0]
			err = sheet.ForEachRow(func(r *xlsx.Row) error {
				var values = []string{}
				err = r.ForEachCell(func(c *xlsx.Cell) error {
					values = append(values, c.Value)
					return nil
				})
				if err != nil {
					return err
				}
				if len(values) == 0 {
					return nil
				}
				if values[0] == "开始IP" || values[0] == "IP" {
					return nil
				}
				item := this.createItemFromValues(values, &countIgnore)
				if item != nil {
					items = append(items, item)
				}
				return nil
			})
			if err != nil {
				this.Fail("Excel读取错误：" + err.Error())
			}
		}
	case ".csv":
		reader := csv.NewReader(bytes.NewBuffer(data))
		for {
			values, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					break
				}
				this.Fail("CSV读取错误：" + err.Error())
			}
			item := this.createItemFromValues(values, &countIgnore)
			if item != nil {
				items = append(items, item)
			}
		}
	case ".json":
		err = json.Unmarshal(data, &items)
		if err != nil {
			this.Fail("导入失败：" + err.Error())
		}
	case ".txt":
		lines := bytes.Split(data, []byte{'\n'})
		for _, line := range lines {
			if len(line) == 0 {
				continue
			}
			item := this.createItemFromValues(strings.SplitN(string(line), ",", 6), &countIgnore)
			if item != nil {
				items = append(items, item)
			}
		}
	}

	var count = 0

	lists.Reverse(items)

	for _, item := range items {
		_, err = this.RPC().IPItemRPC().CreateIPItem(this.AdminContext(), &pb.CreateIPItemRequest{
			IpListId:   params.ListId,
			IpFrom:     item.IpFrom,
			IpTo:       item.IpTo,
			ExpiredAt:  item.ExpiredAt,
			Reason:     item.Reason,
			Type:       item.Type,
			EventLevel: item.EventLevel,
		})
		if err != nil {
			this.Fail("导入过程中出错：" + err.Error())
		}
		count++
	}

	this.Data["count"] = count
	this.Data["countIgnore"] = countIgnore

	this.Success()
}

func (this *ImportAction) createItemFromValues(values []string, countIgnore *int) *pb.IPItem {
	// ipFrom, ipTo, expiredAt, type, eventType, reason

	var item = &pb.IPItem{}
	switch len(values) {
	case 1:
		item.IpFrom = values[0]
	case 2:
		item.IpFrom = values[0]
		item.IpTo = values[1]
	case 3:
		item.IpFrom = values[0]
		item.IpTo = values[1]
		item.ExpiredAt = types.Int64(values[2])
	case 4:
		item.IpFrom = values[0]
		item.IpTo = values[1]
		item.ExpiredAt = types.Int64(values[2])
		item.Type = values[3]
	case 5:
		item.IpFrom = values[0]
		item.IpTo = values[1]
		item.ExpiredAt = types.Int64(values[2])
		item.Type = values[3]
		item.EventLevel = values[4]
	case 6:
		item.IpFrom = values[0]
		item.IpTo = values[1]
		item.ExpiredAt = types.Int64(values[2])
		item.Type = values[3]
		item.EventLevel = values[4]
		item.Reason = values[5]
	}

	if net.ParseIP(item.IpFrom) == nil {
		*countIgnore++
		return nil
	}
	if len(item.IpTo) > 0 && net.ParseIP(item.IpTo) == nil {
		*countIgnore++
		return nil
	}
	return item
}
