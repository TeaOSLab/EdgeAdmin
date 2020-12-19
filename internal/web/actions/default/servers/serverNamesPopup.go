package servers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

type ServerNamesPopupAction struct {
	actionutils.ParentAction
}

func (this *ServerNamesPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *ServerNamesPopupAction) RunGet(params struct {
	ServerId int64
}) {
	serverNamesResp, err := this.RPC().ServerRPC().FindServerNames(this.AdminContext(), &pb.FindServerNamesRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if serverNamesResp.IsAuditing {
		serverNamesResp.ServerNamesJSON = serverNamesResp.AuditingServerNamesJSON
	}
	serverNames := []*serverconfigs.ServerNameConfig{}
	if len(serverNamesResp.ServerNamesJSON) > 0 {
		err = json.Unmarshal(serverNamesResp.ServerNamesJSON, &serverNames)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	result := []string{}
	for _, serverName := range serverNames {
		if len(serverName.SubNames) == 0 {
			result = append(result, serverName.Name)
		} else {
			result = append(result, serverName.SubNames...)
		}
	}
	this.Data["serverNames"] = result

	this.Show()
}
