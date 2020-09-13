package grants

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/grants/grantutils"
	"github.com/iwind/TeaGo/maps"
)

type GrantAction struct {
	actionutils.ParentAction
}

func (this *GrantAction) Init() {
	this.Nav("", "grant", "index")
}

func (this *GrantAction) RunGet(params struct {
	GrantId int64
}) {
	grantResp, err := this.RPC().NodeGrantRPC().FindEnabledGrant(this.AdminContext(), &pb.FindEnabledGrantRequest{GrantId: params.GrantId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if grantResp.Grant == nil {
		this.WriteString("can not find the grant")
		return
	}

	// TODO 处理节点专用的认证

	grant := grantResp.Grant
	this.Data["grant"] = maps.Map{
		"id":          grant.Id,
		"name":        grant.Name,
		"method":      grant.Method,
		"methodName":  grantutils.FindGrantMethodName(grant.Method),
		"username":    grant.Username,
		"password":    grant.Password,
		"privateKey":  grant.PrivateKey,
		"description": grant.Description,
		"su":          grant.Su,
	}

	this.Show()
}
