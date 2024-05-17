// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package webp

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("webp")
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
	this.Data["hasGroupConfig"] = groupResp.HasWebPConfig
	this.Data["groupSettingURL"] = "/servers/groups/group/settings/webp?groupId=" + types.String(groupResp.ServerGroupId)

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["webpConfig"] = webConfig.WebP

	// WebP策略配置
	var serverMap = this.Data.GetMap("server")
	var clusterId = serverMap.GetInt64("clusterId")
	webpPolicyResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterWebPPolicy(this.AdminContext(), &pb.FindEnabledNodeClusterWebPPolicyRequest{NodeClusterId: clusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["requireCache"] = true
	if len(webpPolicyResp.WebpPolicyJSON) > 0 {
		var webpPolicy = &nodeconfigs.WebPImagePolicy{}
		err = json.Unmarshal(webpPolicyResp.WebpPolicyJSON, webpPolicy)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.Data["requireCache"] = webpPolicy.RequireCache
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId    int64
	WebpJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var webpConfig = &serverconfigs.WebPImageConfig{}
	err := json.Unmarshal(params.WebpJSON, webpConfig)
	if err != nil {
		this.Fail("参数校验失败：" + err.Error())
	}

	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebWebP(this.AdminContext(), &pb.UpdateHTTPWebWebPRequest{
		HttpWebId: params.WebId,
		WebpJSON:  params.WebpJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
