package helpers

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	nodes "github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/setup"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"net/http"
	"reflect"
)

// 认证拦截
type UserMustAuth struct {
	AdminId int
	Grant   string
}

func NewUserMustAuth() *UserMustAuth {
	return &UserMustAuth{}
}

func (this *UserMustAuth) BeforeAction(actionPtr actions.ActionWrapper, paramName string) (goNext bool) {
	var action = actionPtr.Object()

	// 安全相关
	if !teaconst.EnabledFrame {
		action.AddHeader("X-Frame-Options", "SAMEORIGIN")
	}
	action.AddHeader("Content-Security-Policy", "default-src 'self' data:; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'")

	// 检查系统是否已经配置过
	if !setup.IsConfigured() {
		action.RedirectURL("/setup")
		return
	}

	var session = action.Session()
	var adminId = session.GetInt("adminId")
	if adminId <= 0 {
		this.login(action)
		return false
	}

	// 检查用户是否存在
	rpc, err := nodes.SharedRPC()
	if err != nil {
		action.WriteString("setup rpc error: " + err.Error())
		utils.PrintError(err)
		return false
	}

	rpcResp, err := rpc.AdminRPC().CheckAdminExists(rpc.Context(0), &pb.CheckAdminExistsRequest{AdminId: int64(adminId)})
	if err != nil {
		utils.PrintError(err)
		action.WriteString(teaconst.ErrServer)
		return false
	}

	if !rpcResp.IsOk {
		session.Delete()

		this.login(action)
		return false
	}

	this.AdminId = adminId
	action.Context.Set("adminId", this.AdminId)

	if action.Request.Method != http.MethodGet {
		return true
	}

	// 初始化内置方法
	action.ViewFunc("teaTitle", func() string {
		return action.Data["teaTitle"].(string)
	})

	// 初始化变量
	modules := []maps.Map{
		{
			"code": "servers",
			"name": "代理服务",
			"icon": "clone outsize",
			"subItems": []maps.Map{
				{
					"name": "通用组件",
					"url":  "/servers/components",
					"code": "components",
				},
			},
		},
		{
			"code": "clusters",
			"name": "节点集群",
			"icon": "cloud",
		},
		{
			"code": "dns",
			"name": "DNS",
			"icon": "globe",
		},
		{
			"code": "settings",
			"name": "系统设置",
			"icon": "setting",
		},
	}

	action.Data["teaTitle"] = teaconst.ProductNameZH
	action.Data["teaName"] = teaconst.ProductNameZH

	resp, err := rpc.AdminRPC().FindAdminFullname(rpc.Context(0), &pb.FindAdminFullnameRequest{AdminId: int64(this.AdminId)})
	if err != nil {
		utils.PrintError(err)
		action.Data["teaUsername"] = ""
	} else {
		action.Data["teaUsername"] = resp.Fullname
	}

	action.Data["teaUserAvatar"] = ""

	action.Data["teaMenu"] = ""
	action.Data["teaModules"] = modules
	action.Data["teaSubMenus"] = []map[string]interface{}{}
	action.Data["teaTabbar"] = []map[string]interface{}{}
	action.Data["teaVersion"] = teaconst.Version
	action.Data["teaIsSuper"] = false
	action.Data["teaDemoEnabled"] = teaconst.IsDemo
	action.Data["teaSubMenu"] = ""

	// 菜单
	action.Data["firstMenuItem"] = ""

	// 未读消息数
	action.Data["teaBadge"] = 0

	// 调用Init
	initMethod := reflect.ValueOf(actionPtr).MethodByName("Init")
	if initMethod.IsValid() {
		initMethod.Call([]reflect.Value{})
	}

	return true
}

func (this *UserMustAuth) login(action *actions.ActionObject) {
	action.RedirectURL("/")
}
