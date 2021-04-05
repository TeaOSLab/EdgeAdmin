package recipients

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

// 媒介类型选项
type MediaOptionsAction struct {
	actionutils.ParentAction
}

func (this *MediaOptionsAction) RunPost(params struct{}) {
	resp, err := this.RPC().MessageMediaRPC().FindAllMessageMedias(this.AdminContext(), &pb.FindAllMessageMediasRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	mediaMaps := []maps.Map{}
	for _, media := range resp.MessageMedias {
		mediaMaps = append(mediaMaps, maps.Map{
			"id":          media.Id,
			"type":        media.Type,
			"name":        media.Name,
			"description": media.Description,
		})
	}
	this.Data["medias"] = mediaMaps

	this.Success()
}
