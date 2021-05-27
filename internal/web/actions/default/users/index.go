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
	Keyword string
}) {
	countResp, err := this.RPC().UserRPC().CountAllEnabledUsers(this.AdminContext(), &pb.CountAllEnabledUsersRequest{Keyword: params.Keyword})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	usersResp, err := this.RPC().UserRPC().ListEnabledUsers(this.AdminContext(), &pb.ListEnabledUsersRequest{
		Keyword: params.Keyword,
		Offset:  page.Offset,
		Size:    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	userMaps := []maps.Map{}
	for _, user := range usersResp.Users {
		var clusterMap maps.Map = nil
		if user.NodeCluster != nil {
			clusterMap = maps.Map{
				"id":   user.NodeCluster.Id,
				"name": user.NodeCluster.Name,
			}
		}
		userMaps = append(userMaps, maps.Map{
			"id":          user.Id,
			"username":    user.Username,
			"isOn":        user.IsOn,
			"fullname":    user.Fullname,
			"email":       user.Email,
			"mobile":      user.Mobile,
			"tel":         user.Tel,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", user.CreatedAt),
			"cluster":     clusterMap,
		})
	}
	this.Data["users"] = userMaps

	this.Show()
}
