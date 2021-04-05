package instances

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

// 媒介类型选项
type OptionsAction struct {
	actionutils.ParentAction
}

func (this *OptionsAction) RunPost(params struct{}) {
	resp, err := this.RPC().MessageMediaInstanceRPC().ListEnabledMessageMediaInstances(this.AdminContext(), &pb.ListEnabledMessageMediaInstancesRequest{
		Offset: 0,
		Size:   1000,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	instanceMaps := []maps.Map{}
	for _, instance := range resp.MessageMediaInstances {
		instanceMaps = append(instanceMaps, maps.Map{
			"id":          instance.Id,
			"name":        instance.Name,
			"description": instance.Description,
			"media": maps.Map{
				"type":            instance.MessageMedia.Type,
				"name":            instance.MessageMedia.Name,
				"userDescription": instance.MessageMedia.UserDescription,
			},
		})
	}
	this.Data["instances"] = instanceMaps

	this.Success()
}
