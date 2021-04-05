package admins

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

// 系统用户选项
// 组件中需要用到的系统用户选项
type OptionsAction struct {
	actionutils.ParentAction
}

func (this *OptionsAction) Init() {
	this.Nav("", "", "")
}

func (this *OptionsAction) RunPost(params struct{}) {
	// TODO 实现关键词搜索
	adminsResp, err := this.RPC().AdminRPC().ListEnabledAdmins(this.AdminContext(), &pb.ListEnabledAdminsRequest{
		Offset: 0,
		Size:   1000,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	adminMaps := []maps.Map{}
	for _, admin := range adminsResp.Admins {
		adminMaps = append(adminMaps, maps.Map{
			"id":       admin.Id,
			"name":     admin.Fullname,
			"username": admin.Username,
		})
	}
	this.Data["admins"] = adminMaps

	this.Success()
}
