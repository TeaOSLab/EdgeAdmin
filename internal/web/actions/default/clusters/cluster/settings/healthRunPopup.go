package settings

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type HealthRunPopupAction struct {
	actionutils.ParentAction
}

func (this *HealthRunPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *HealthRunPopupAction) RunGet(params struct{}) {

	this.Show()
}

func (this *HealthRunPopupAction) RunPost(params struct {
	ClusterId int64

	Must *actions.Must
}) {
	resp, err := this.RPC().NodeClusterRPC().ExecuteNodeClusterHealthCheck(this.AdminContext(), &pb.ExecuteNodeClusterHealthCheckRequest{ClusterId: params.ClusterId})
	if err != nil {
		this.Fail(err.Error())
	}

	this.Data["results"] = resp.Results
	this.Success()
}
