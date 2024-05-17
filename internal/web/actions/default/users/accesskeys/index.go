// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package accesskeys

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/users/userutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "accessKey")
}

func (this *IndexAction) RunGet(params struct {
	UserId int64
}) {
	err := userutils.InitUser(this.Parent(), params.UserId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	accessKeysResp, err := this.RPC().UserAccessKeyRPC().FindAllEnabledUserAccessKeys(this.AdminContext(), &pb.FindAllEnabledUserAccessKeysRequest{UserId: params.UserId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	accessKeyMaps := []maps.Map{}
	for _, accessKey := range accessKeysResp.UserAccessKeys {
		var accessedTime string
		if accessKey.AccessedAt > 0 {
			accessedTime = timeutil.FormatTime("Y-m-d H:i:s", accessKey.AccessedAt)
		}
		accessKeyMaps = append(accessKeyMaps, maps.Map{
			"id":           accessKey.Id,
			"isOn":         accessKey.IsOn,
			"uniqueId":     accessKey.UniqueId,
			"secret":       accessKey.Secret,
			"description":  accessKey.Description,
			"accessedTime": accessedTime,
		})
	}
	this.Data["accessKeys"] = accessKeyMaps

	this.Show()
}
