package servers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "server", "index")
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
		config := &serverconfigs.ServerConfig{}
		err = json.Unmarshal(server.Config, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 端口列表
		portMaps := []maps.Map{}
		if len(server.HttpJSON) > 0 && config.HTTP.IsOn {
			for _, listen := range config.HTTP.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if config.HTTPS != nil && config.HTTPS.IsOn {
			for _, listen := range config.HTTPS.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if config.TCP != nil && config.TCP.IsOn {
			for _, listen := range config.TCP.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if config.TLS != nil && config.TLS.IsOn {
			for _, listen := range config.TLS.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if config.Unix != nil && config.Unix.IsOn {
			for _, listen := range config.Unix.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.Host,
				})
			}
		}
		if config.UDP != nil && config.UDP.IsOn {
			for _, listen := range config.UDP.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}

		serverMaps = append(serverMaps, maps.Map{
			"id":   server.Id,
			"isOn": server.IsOn,
			"name": server.Name,
			"cluster": maps.Map{
				"id":   server.Cluster.Id,
				"name": server.Cluster.Name,
			},
			"ports":          portMaps,
			"serverTypeName": serverconfigs.FindServerType(server.Type).GetString("name"),
		})
	}
	this.Data["servers"] = serverMaps

	this.Show()
}
