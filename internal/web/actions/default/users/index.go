package users

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
	Keyword          string
	Verifying        bool
	MobileIsVerified int32 `default:"-1"`
}) {
	this.Data["keyword"] = params.Keyword
	this.Data["isVerifying"] = params.Verifying
	this.Data["mobileIsVerified"] = params.MobileIsVerified

	// 未审核的总数量
	countVerifyingUsersResp, err := this.RPC().UserRPC().CountAllEnabledUsers(this.AdminContext(), &pb.CountAllEnabledUsersRequest{
		IsVerifying: true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["countVerifyingUsers"] = countVerifyingUsersResp.Count

	// 当前匹配的数量
	countResp, err := this.RPC().UserRPC().CountAllEnabledUsers(this.AdminContext(), &pb.CountAllEnabledUsersRequest{
		Keyword:          params.Keyword,
		IsVerifying:      params.Verifying,
		MobileIsVerified: params.MobileIsVerified,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	usersResp, err := this.RPC().UserRPC().ListEnabledUsers(this.AdminContext(), &pb.ListEnabledUsersRequest{
		Keyword:          params.Keyword,
		IsVerifying:      params.Verifying,
		MobileIsVerified: params.MobileIsVerified,
		Offset:           page.Offset,
		Size:             page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var userMaps = []maps.Map{}
	for _, user := range usersResp.Users {
		var clusterMap maps.Map = nil
		if user.NodeCluster != nil {
			clusterMap = maps.Map{
				"id":   user.NodeCluster.Id,
				"name": user.NodeCluster.Name,
			}
		}

		isSubmittedResp, err := this.RPC().UserIdentityRPC().CheckUserIdentityIsSubmitted(this.AdminContext(), &pb.CheckUserIdentityIsSubmittedRequest{UserId: user.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var identityIsSubmitted = isSubmittedResp.IsSubmitted

		userMaps = append(userMaps, maps.Map{
			"id":                  user.Id,
			"username":            user.Username,
			"isOn":                user.IsOn,
			"fullname":            user.Fullname,
			"email":               user.Email,
			"mobile":              user.Mobile,
			"tel":                 user.Tel,
			"createdTime":         timeutil.FormatTime("Y-m-d H:i:s", user.CreatedAt),
			"cluster":             clusterMap,
			"registeredIP":        user.RegisteredIP,
			"isVerified":          user.IsVerified,
			"isRejected":          user.IsRejected,
			"identityIsSubmitted": identityIsSubmitted,
			"verifiedMobile":      user.VerifiedMobile,
		})
	}
	this.Data["users"] = userMaps

	this.Show()
}
