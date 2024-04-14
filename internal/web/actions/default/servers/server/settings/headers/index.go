package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/types"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("header")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// 只有HTTP服务才支持
	if this.FilterHTTPFamily() {
		return
	}

	// 服务分组设置
	groupResp, err := this.RPC().ServerGroupRPC().FindEnabledServerGroupConfigInfo(this.AdminContext(), &pb.FindEnabledServerGroupConfigInfoRequest{
		ServerId: params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasGroupRequestConfig"] = groupResp.HasRequestHeadersConfig
	this.Data["hasGroupResponseConfig"] = groupResp.HasResponseHeadersConfig
	this.Data["groupSettingURL"] = "/servers/groups/group/settings/headers?groupId=" + types.String(groupResp.ServerGroupId)

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	webId := webConfig.Id

	isChanged := false
	if webConfig.RequestHeaderPolicy == nil {
		createHeaderPolicyResp, err := this.RPC().HTTPHeaderPolicyRPC().CreateHTTPHeaderPolicy(this.AdminContext(), &pb.CreateHTTPHeaderPolicyRequest{})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var headerPolicyId = createHeaderPolicyResp.HttpHeaderPolicyId
		ref := &shared.HTTPHeaderPolicyRef{
			IsPrior:        false,
			IsOn:           true,
			HeaderPolicyId: headerPolicyId,
		}
		refJSON, err := json.Marshal(ref)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebRequestHeader(this.AdminContext(), &pb.UpdateHTTPWebRequestHeaderRequest{
			HttpWebId:  webId,
			HeaderJSON: refJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		isChanged = true
	}
	if webConfig.ResponseHeaderPolicy == nil {
		createHeaderPolicyResp, err := this.RPC().HTTPHeaderPolicyRPC().CreateHTTPHeaderPolicy(this.AdminContext(), &pb.CreateHTTPHeaderPolicyRequest{})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		headerPolicyId := createHeaderPolicyResp.HttpHeaderPolicyId
		ref := &shared.HTTPHeaderPolicyRef{
			IsPrior:        false,
			IsOn:           true,
			HeaderPolicyId: headerPolicyId,
		}
		refJSON, err := json.Marshal(ref)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebResponseHeader(this.AdminContext(), &pb.UpdateHTTPWebResponseHeaderRequest{
			HttpWebId:  webId,
			HeaderJSON: refJSON,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		isChanged = true
	}

	// 重新获取配置
	if isChanged {
		webConfig, err = dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["requestHeaderRef"] = webConfig.RequestHeaderPolicyRef
	this.Data["requestHeaderPolicy"] = webConfig.RequestHeaderPolicy
	this.Data["responseHeaderRef"] = webConfig.ResponseHeaderPolicyRef
	this.Data["responseHeaderPolicy"] = webConfig.ResponseHeaderPolicy

	this.Show()
}
