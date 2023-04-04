package admins

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct {
	Keyword         string
	HasWeakPassword bool
}) {
	this.Data["keyword"] = params.Keyword
	this.Data["hasWeakPassword"] = params.HasWeakPassword

	countResp, err := this.RPC().AdminRPC().CountAllEnabledAdmins(this.AdminContext(), &pb.CountAllEnabledAdminsRequest{
		Keyword:         params.Keyword,
		HasWeakPassword: params.HasWeakPassword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var page = this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	adminsResp, err := this.RPC().AdminRPC().ListEnabledAdmins(this.AdminContext(), &pb.ListEnabledAdminsRequest{
		Keyword:         params.Keyword,
		HasWeakPassword: params.HasWeakPassword,
		Offset:          page.Offset,
		Size:            page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var adminMaps = []maps.Map{}
	for _, admin := range adminsResp.Admins {
		adminMaps = append(adminMaps, maps.Map{
			"id":              admin.Id,
			"isOn":            admin.IsOn,
			"isSuper":         admin.IsSuper,
			"username":        admin.Username,
			"fullname":        admin.Fullname,
			"createdTime":     timeutil.FormatTime("Y-m-d H:i:s", admin.CreatedAt),
			"otpLoginIsOn":    admin.OtpLogin != nil && admin.OtpLogin.IsOn,
			"canLogin":        admin.CanLogin,
			"hasWeakPassword": admin.HasWeakPassword,
		})
	}
	this.Data["admins"] = adminMaps

	this.Show()
}
