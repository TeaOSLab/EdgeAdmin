package helpers

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	nodes "github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
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

	var session = action.Session()
	var adminId = session.GetInt("adminId")
	if adminId <= 0 {
		this.login(action)
		return false
	}

	// 检查用户是否存在
	rpc, err := nodes.SharedRPC()
	if err != nil {
		utils.PrintError(err)
		return false
	}

	rpcResp, err := rpc.AdminRPC().CheckAdminExists(rpc.Context(0), &pb.CheckAdminExistsRequest{AdminId: int64(adminId)})
	if err != nil {
		utils.PrintError(err)
		actionPtr.Object().WriteString(teaconst.ErrServer)
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
	modules := []map[string]interface{}{
		{
			"code":     "servers",
			"menuName": "代理服务",
			"icon":     "clone outsize",
		},
		{
			"code":     "clusters",
			"menuName": "节点集群",
			"icon":     "cloud",
		},
		{
			"code":     "api",
			"menuName": "API节点",
			"icon":     "exchange",
		},
		{
			"code":     "db",
			"menuName": "数据库节点",
			"icon":     "database",
		},
		{
			"code":     "log",
			"menuName": "日志节点",
			"icon":     "dot circle",
		},
		{
			"code":     "dns",
			"menuName": "DNS",
			"icon":     "globe",
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
