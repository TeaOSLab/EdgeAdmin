package node

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"strconv"
)

type Helper struct {
}

func NewHelper() *Helper {
	return &Helper{}
}

func (this *Helper) BeforeAction(action *actions.ActionObject) (goNext bool) {
	if action.Request.Method != http.MethodGet {
		return true
	}

	action.Data["teaMenu"] = "api"

	nodeId := action.ParamInt64("nodeId")
	nodeIdString := strconv.FormatInt(nodeId, 10)

	// 节点信息
	rpcClient, err := rpc.SharedRPC()
	if err != nil {
		logs.Error(err)
		return
	}
	nodeResp, err := rpcClient.APINodeRPC().FindEnabledAPINode(rpcClient.Context(action.Context.GetInt64("adminId")), &pb.FindEnabledAPINodeRequest{NodeId: nodeId})
	if err != nil {
		action.WriteString(err.Error())
		return
	}
	if nodeResp.Node == nil {
		action.WriteString("node not found")
		return
	}

	// 左侧菜单栏
	secondMenuItem := action.Data.GetString("secondMenuItem")
	switch action.Data.GetString("firstMenuItem") {
	case "setting":
		action.Data["leftMenuItems"] = this.createSettingMenu(nodeIdString, secondMenuItem)
	}

	return true
}

// 设置相关菜单
func (this *Helper) createSettingMenu(nodeIdString string, selectedItem string) (items []maps.Map) {
	items = append(items, maps.Map{
		"name":     "基础设置",
		"url":      "/api/node/settings?nodeId=" + nodeIdString,
		"isActive": selectedItem == "basic",
	})
	return
}
