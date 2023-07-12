// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package remoteAddr

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/groups/group/servergrouputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"strings"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("remoteAddr")
}

func (this *IndexAction) RunGet(params struct {
	GroupId int64
}) {
	_, err := servergrouputils.InitGroup(this.Parent(), params.GroupId, "remoteAddr")
	if err != nil {
		this.ErrorPage(err)
		return
	}

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerGroupId(this.AdminContext(), params.GroupId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["remoteAddrConfig"] = webConfig.RemoteAddr

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId          int64
	RemoteAddrJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var remoteAddrConfig = &serverconfigs.HTTPRemoteAddrConfig{}
	err := json.Unmarshal(params.RemoteAddrJSON, remoteAddrConfig)
	if err != nil {
		this.Fail("参数校验失败：" + err.Error())
	}

	remoteAddrConfig.Value = strings.TrimSpace(remoteAddrConfig.Value)
	err = remoteAddrConfig.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
		return
	}

	remoteAddrJSON, err := json.Marshal(remoteAddrConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebRemoteAddr(this.AdminContext(), &pb.UpdateHTTPWebRemoteAddrRequest{
		HttpWebId:      params.WebId,
		RemoteAddrJSON: remoteAddrJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
