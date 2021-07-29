// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipbox

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().HTTPAccessLogPolicyRPC().CountAllEnabledHTTPAccessLogPolicies(this.AdminContext(), &pb.CountAllEnabledHTTPAccessLogPoliciesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	page := this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	policiesResp, err := this.RPC().HTTPAccessLogPolicyRPC().ListEnabledHTTPAccessLogPolicies(this.AdminContext(), &pb.ListEnabledHTTPAccessLogPoliciesRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	var policyMaps = []maps.Map{}
	for _, policy := range policiesResp.HttpAccessLogPolicies {
		var optionsMap = maps.Map{}
		if len(policy.OptionsJSON) > 0 {
			err = json.Unmarshal(policy.OptionsJSON, &optionsMap)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}
		policyMaps = append(policyMaps, maps.Map{
			"id":       policy.Id,
			"name":     policy.Name,
			"type":     policy.Type,
			"typeName": serverconfigs.FindAccessLogStorageTypeName(policy.Type),
			"isOn":     policy.IsOn,
			"isPublic": policy.IsPublic,
			"options":  optionsMap,
		})
	}
	this.Data["policies"] = policyMaps

	this.Show()
}
