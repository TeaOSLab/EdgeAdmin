package webutils

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

// 根据ServerId查找Web配置
func FindWebConfigWithServerId(parentAction *actionutils.ParentAction, serverId int64) (*serverconfigs.HTTPWebConfig, error) {
	resp, err := parentAction.RPC().ServerRPC().FindAndInitServerWebConfig(parentAction.AdminContext(), &pb.FindAndInitServerWebConfigRequest{ServerId: serverId})
	if err != nil {
		return nil, err
	}
	config := &serverconfigs.HTTPWebConfig{}
	err = json.Unmarshal(resp.WebJSON, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// 根据LocationId查找Web配置
func FindWebConfigWithLocationId(parentAction *actionutils.ParentAction, locationId int64) (*serverconfigs.HTTPWebConfig, error) {
	resp, err := parentAction.RPC().HTTPLocationRPC().FindAndInitHTTPLocationWebConfig(parentAction.AdminContext(), &pb.FindAndInitHTTPLocationWebConfigRequest{LocationId: locationId})
	if err != nil {
		return nil, err
	}
	config := &serverconfigs.HTTPWebConfig{}
	err = json.Unmarshal(resp.WebJSON, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

// 根据WebId查找Web配置
func FindWebConfigWithId(parentAction *actionutils.ParentAction, webId int64) (*serverconfigs.HTTPWebConfig, error) {
	resp, err := parentAction.RPC().HTTPWebRPC().FindEnabledHTTPWebConfig(parentAction.AdminContext(), &pb.FindEnabledHTTPWebConfigRequest{WebId: webId})
	if err != nil {
		return nil, err
	}
	config := &serverconfigs.HTTPWebConfig{}
	err = json.Unmarshal(resp.WebJSON, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
