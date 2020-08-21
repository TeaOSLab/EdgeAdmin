package servers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs/serverconfigs"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc/pb"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
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
		// 服务名
		serverConfig := &serverconfigs.ServerConfig{}
		err = json.Unmarshal(server.Config, &serverConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		err = serverConfig.Init()
		if err != nil {
			logs.Println("init server '" + serverConfig.Name + "' error: " + err.Error())
		}

		serverTypeNames := []string{}

		// 端口列表
		portMaps := []maps.Map{}
		if serverConfig.HTTP != nil && serverConfig.HTTP.IsOn {
			serverTypeNames = append(serverTypeNames, "HTTP")
			for _, listen := range serverConfig.HTTP.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if serverConfig.HTTPS != nil && serverConfig.HTTPS.IsOn {
			serverTypeNames = append(serverTypeNames, "HTTPS")
			for _, listen := range serverConfig.HTTPS.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if serverConfig.TCP != nil && serverConfig.TCP.IsOn {
			serverTypeNames = append(serverTypeNames, "TCP")
			for _, listen := range serverConfig.TCP.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if serverConfig.TLS != nil && serverConfig.TLS.IsOn {
			serverTypeNames = append(serverTypeNames, "TLS")
			for _, listen := range serverConfig.TLS.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}
		if serverConfig.Unix != nil && serverConfig.Unix.IsOn {
			serverTypeNames = append(serverTypeNames, "Unix")
			for _, listen := range serverConfig.Unix.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.Host,
				})
			}
		}
		if serverConfig.UDP != nil && serverConfig.UDP.IsOn {
			serverTypeNames = append(serverTypeNames, "UDP")
			for _, listen := range serverConfig.UDP.Listen {
				portMaps = append(portMaps, maps.Map{
					"protocol":  listen.Protocol,
					"portRange": listen.PortRange,
				})
			}
		}

		serverMaps = append(serverMaps, maps.Map{
			"id":   server.Id,
			"name": serverConfig.Name,
			"cluster": maps.Map{
				"id":   server.Cluster.Id,
				"name": server.Cluster.Name,
			},
			"ports":          portMaps,
			"serverTypeName": strings.Join(serverTypeNames, "+"),
		})
	}
	this.Data["servers"] = serverMaps

	this.Show()
}
