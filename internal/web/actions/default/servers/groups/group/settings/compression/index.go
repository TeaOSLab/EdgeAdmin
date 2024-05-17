// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package compression

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/groups/group/servergrouputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("compression")
}

func (this *IndexAction) RunGet(params struct {
	GroupId int64
}) {
	_, err := servergrouputils.InitGroup(this.Parent(), params.GroupId, "compression")
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// WebId
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerGroupId(this.AdminContext(), params.GroupId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id

	this.Data["compressionConfig"] = webConfig.Compression

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId           int64
	CompressionJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.ServerCompression_LogUpdateCompressionSettings, params.WebId)

	// 校验配置
	var compressionConfig = &serverconfigs.HTTPCompressionConfig{}
	err := json.Unmarshal(params.CompressionJSON, compressionConfig)
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
	}

	err = compressionConfig.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
	}

	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebCompression(this.AdminContext(), &pb.UpdateHTTPWebCompressionRequest{
		HttpWebId:       params.WebId,
		CompressionJSON: params.CompressionJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
