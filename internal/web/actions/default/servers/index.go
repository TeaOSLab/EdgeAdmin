package servers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs/nodes"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "index")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().ServerRPC().CountAllEnabledServers(this.AdminContext(), &pb.CountAllEnabledServersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	// 服务列表
	serversResp, err := this.RPC().ServerRPC().ListEnabledServers(this.AdminContext(), &pb.ListEnabledServersRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	serverMaps := []maps.Map{}
	for _, server := range serversResp.Servers {
		// 服务名
		serverConfig := &nodes.ServerConfig{}
		err = json.Unmarshal(server.Config, &serverConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		serverMaps = append(serverMaps, maps.Map{
			"id":   server.Id,
			"name": serverConfig.Name,
			"cluster": maps.Map{
				"id":   server.Cluster.Id,
				"name": server.Cluster.Name,
			},
		})
	}
	this.Data["servers"] = serverMaps

	this.Show()
}
