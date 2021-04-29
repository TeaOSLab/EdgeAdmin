// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type MemoryAction struct {
	actionutils.ParentAction
}

func (this *MemoryAction) RunPost(params struct {
	NodeId int64
}) {
	resp, err := this.RPC().NodeValueRPC().ListNodeValues(this.AdminContext(), &pb.ListNodeValuesRequest{
		Role:   "node",
		NodeId: params.NodeId,
		Item:   nodeconfigs.NodeValueItemMemory,
		Range:  nodeconfigs.NodeValueRangeMinute,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	valuesMap := map[string]float32{} // YmdHi => usage
	for _, v := range resp.NodeValues {
		if len(v.ValueJSON) == 0 {
			continue
		}

		valueMap := maps.Map{}
		err = json.Unmarshal(v.ValueJSON, &valueMap)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		valuesMap[timeutil.FormatTime("YmdHi", v.CreatedAt)] = valueMap.GetFloat32("usage") * 100
	}

	// 过去一个小时
	result := []maps.Map{}
	for i := 60; i >= 1; i-- {
		timestamp := time.Now().Unix() - int64(i)*60
		minute := timeutil.FormatTime("YmdHi", timestamp)
		total, ok := valuesMap[minute]
		if ok {
			result = append(result, maps.Map{
				"label": timeutil.FormatTime("H:i", timestamp),
				"value": total,
				"text":  fmt.Sprintf("%.2f%%", total),
			})
		} else {
			result = append(result, maps.Map{
				"label": timeutil.FormatTime("H:i", timestamp),
				"value": 0,
				"text":  "0.0%",
			})
		}
	}

	this.Data["values"] = result

	this.Success()
}
