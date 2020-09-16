package reverseProxy

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

type UpdateOnAction struct {
	actionutils.ParentAction
}

func (this *UpdateOnAction) RunPost(params struct {
	ServerId       int64
	ReverseProxyId int64
	IsOn           bool
}) {
	// 如果没有配置过，则配置
	if params.ReverseProxyId <= 0 {
		if !params.IsOn {
			this.Success()
		}

		resp, err := this.RPC().ReverseProxyRPC().CreateReverseProxy(this.AdminContext(), &pb.CreateReverseProxyRequest{
			SchedulingJSON:     nil,
			PrimaryOriginsJSON: nil,
			BackupOriginsJSON:  nil,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		reverseProxyId := resp.ReverseProxyId
		_, err = this.RPC().ServerRPC().UpdateServerReverseProxy(this.AdminContext(), &pb.UpdateServerReverseProxyRequest{
			ServerId:       params.ServerId,
			ReverseProxyId: reverseProxyId,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.Success()
	}

	// 如果已经配置过
	_, err := this.RPC().ReverseProxyRPC().UpdateReverseProxyIsOn(this.AdminContext(), &pb.UpdateReverseProxyIsOnRequest{
		ReverseProxyId: params.ReverseProxyId,
		IsOn:           params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
