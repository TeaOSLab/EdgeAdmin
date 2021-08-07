package instances

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "instance")
}

func (this *IndexAction) RunGet(params struct {
}) {
	// TODO 增加系统用户、媒介类型等条件搜索
	countResp, err := this.RPC().MessageMediaInstanceRPC().CountAllEnabledMessageMediaInstances(this.AdminContext(), &pb.CountAllEnabledMessageMediaInstancesRequest{
		MediaType: "",
		Keyword:   "",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	instancesResp, err := this.RPC().MessageMediaInstanceRPC().ListEnabledMessageMediaInstances(this.AdminContext(), &pb.ListEnabledMessageMediaInstancesRequest{
		MediaType: "",
		Keyword:   "",
		Offset:    page.Offset,
		Size:      page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	instanceMaps := []maps.Map{}
	for _, instance := range instancesResp.MessageMediaInstances {
		if instance.MessageMedia == nil {
			continue
		}
		instanceMaps = append(instanceMaps, maps.Map{
			"id":   instance.Id,
			"name": instance.Name,
			"isOn": instance.IsOn,
			"media": maps.Map{
				"name": instance.MessageMedia.Name,
			},
			"description": instance.Description,
		})
	}
	this.Data["instances"] = instanceMaps

	this.Show()
}
