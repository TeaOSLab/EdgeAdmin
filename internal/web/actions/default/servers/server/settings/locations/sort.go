package locations

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

type SortAction struct {
	actionutils.ParentAction
}

func (this *SortAction) RunPost(params struct {
	WebId       int64
	LocationIds []int64
}) {
	if len(params.LocationIds) == 0 {
		this.Success()
	}

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithId(this.AdminContext(), params.WebId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if webConfig == nil {
		this.Success()
		return
	}
	refMap := map[int64]*serverconfigs.HTTPLocationRef{}
	for _, ref := range webConfig.LocationRefs {
		refMap[ref.LocationId] = ref
	}

	newRefs := []*serverconfigs.HTTPLocationRef{}
	for _, locationId := range params.LocationIds {
		ref, ok := refMap[locationId]
		if ok {
			newRefs = append(newRefs, ref)
		}
	}
	newRefsJSON, err := json.Marshal(newRefs)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebLocations(this.AdminContext(), &pb.UpdateHTTPWebLocationsRequest{
		HttpWebId:     params.WebId,
		LocationsJSON: newRefsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
