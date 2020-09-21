package locations

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/webutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	WebId      int64
	LocationId int64
}) {
	webConfig, err := webutils.FindWebConfigWithId(this.Parent(), params.WebId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	webConfig.RemoveLocationRef(params.LocationId)

	refJSON, err := json.Marshal(webConfig.LocationRefs)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebLocations(this.AdminContext(), &pb.UpdateHTTPWebLocationsRequest{
		WebId:         params.WebId,
		LocationsJSON: refJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
