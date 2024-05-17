// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package webp

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {

}

func (this *IndexAction) RunGet(params struct {
	LocationId int64
}) {
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithLocationId(this.AdminContext(), params.LocationId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["webpConfig"] = webConfig.WebP

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
