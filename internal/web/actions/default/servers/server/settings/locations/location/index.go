package location

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

// 路径规则详情
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.TinyMenu("basic")
}

func (this *IndexAction) RunGet(params struct {
	LocationId int64
}) {
	locationConfigResp, err := this.RPC().HTTPLocationRPC().FindEnabledHTTPLocationConfig(this.AdminContext(), &pb.FindEnabledHTTPLocationConfigRequest{LocationId: params.LocationId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	locationConfig := &serverconfigs.HTTPLocationConfig{}
	err = json.Unmarshal(locationConfigResp.LocationJSON, locationConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["locationId"] = locationConfig.Id
	this.Data["locationConfig"] = locationConfig

	this.Show()
}
