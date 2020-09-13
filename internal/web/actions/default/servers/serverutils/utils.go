package serverutils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"strconv"
)

// 查找Server
func FindServer(p *actionutils.ParentAction, serverId int64) (*pb.Server, *serverconfigs.ServerConfig, bool) {
	serverResp, err := p.RPC().ServerRPC().FindEnabledServer(p.AdminContext(), &pb.FindEnabledServerRequest{ServerId: serverId})
	if err != nil {
		p.ErrorPage(err)
		return nil, nil, false
	}
	server := serverResp.Server
	if server == nil {
		p.ErrorPage(errors.New("not found server with id '" + strconv.FormatInt(serverId, 10) + "'"))
		return nil, nil, false
	}
	config, err := serverconfigs.NewServerConfigFromJSON(server.Config)
	if err != nil {
		p.ErrorPage(err)
		return nil, nil, false
	}

	return server, config, true
}
