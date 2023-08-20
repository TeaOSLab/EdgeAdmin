package httpReverseProxy

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/groups/group/servergrouputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/schedulingconfigs"
)

type SchedulingAction struct {
	actionutils.ParentAction
}

func (this *SchedulingAction) Init() {
	this.FirstMenu("scheduling")
}

func (this *SchedulingAction) RunGet(params struct {
	GroupId int64
}) {
	_, err := servergrouputils.InitGroup(this.Parent(), params.GroupId, "httpReverseProxy")
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["family"] = "http"

	reverseProxyResp, err := this.RPC().ServerGroupRPC().FindAndInitServerGroupHTTPReverseProxyConfig(this.AdminContext(), &pb.FindAndInitServerGroupHTTPReverseProxyConfigRequest{ServerGroupId: params.GroupId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	reverseProxy := serverconfigs.NewReverseProxyConfig()
	err = json.Unmarshal(reverseProxyResp.ReverseProxyJSON, reverseProxy)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["reverseProxyId"] = reverseProxy.Id

	schedulingCode := reverseProxy.FindSchedulingConfig().Code
	schedulingMap := schedulingconfigs.FindSchedulingType(schedulingCode)
	if schedulingMap == nil {
		this.ErrorPage(errors.New("invalid scheduling code '" + schedulingCode + "'"))
		return
	}
	this.Data["scheduling"] = schedulingMap

	this.Show()
}
