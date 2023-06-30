// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package updates

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/systemconfigs"
)

type ResetIgnoredVersionAction struct {
	actionutils.ParentAction
}

func (this *ResetIgnoredVersionAction) RunPost(params struct{}) {
	defer this.CreateLogInfo(codes.AdminUpdate_LogResetIgnoreVersion)

	valueResp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: systemconfigs.SettingCodeCheckUpdates})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var valueJSON = valueResp.ValueJSON
	var config = systemconfigs.NewCheckUpdatesConfig()
	if len(valueJSON) > 0 {
		err = json.Unmarshal(valueJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	config.IgnoredVersion = ""
	configJSON, err := json.Marshal(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().SysSettingRPC().UpdateSysSetting(this.AdminContext(), &pb.UpdateSysSettingRequest{
		Code:      systemconfigs.SettingCodeCheckUpdates,
		ValueJSON: configJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
