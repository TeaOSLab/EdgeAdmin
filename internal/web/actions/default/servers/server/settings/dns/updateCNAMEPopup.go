// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package dns

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/types"
	"regexp"
	"strings"
)

type UpdateCNAMEPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateCNAMEPopupAction) RunGet(params struct {
	ServerId int64
}) {
	this.Data["serverId"] = params.ServerId

	dnsInfoResp, err := this.RPC().ServerRPC().FindEnabledServerDNS(this.AdminContext(), &pb.FindEnabledServerDNSRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["dnsName"] = dnsInfoResp.DnsName

	this.Show()
}

func (this *UpdateCNAMEPopupAction) RunPost(params struct {
	ServerId int64
	DnsName  string
}) {
	defer this.CreateLogInfo(codes.ServerDNS_LogUpdateDNSName, params.ServerId, params.DnsName)

	var dnsName = strings.ToLower(params.DnsName)
	if len(dnsName) == 0 {
		this.FailField("dnsName", "CNAME不能为空")
	}

	const maxLen = 30
	if len(dnsName) > maxLen {
		this.FailField("dnsName", "CNAME长度不能超过"+types.String(maxLen)+"个字符")
	}
	if !regexp.MustCompile(`^[a-z0-9]{1,` + types.String(maxLen) + `}$`).MatchString(dnsName) {
		this.FailField("dnsName", "CNAME中只能包含数字、英文字母")
	}

	serverResp, err := this.RPC().ServerRPC().FindEnabledServer(this.AdminContext(), &pb.FindEnabledServerRequest{
		ServerId:       params.ServerId,
		IgnoreSSLCerts: true,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var server = serverResp.Server
	if server == nil {
		this.Fail("找不到要修改的服务")
	}
	if server.NodeCluster == nil {
		this.Fail("服务必须先分配到一个集群才能修改")
	}
	var clusterId = server.NodeCluster.Id

	if server.DnsName == params.DnsName {
		// 没有修改则直接返回
		this.Success()
	}

	serverIdResp, err := this.RPC().ServerRPC().FindServerIdWithDNSName(this.AdminContext(), &pb.FindServerIdWithDNSNameRequest{
		NodeClusterId: clusterId,
		DnsName:       dnsName,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	if serverIdResp.ServerId > 0 && serverIdResp.ServerId != params.ServerId {
		this.FailField("dnsName", "当前CNAME已被别的服务占用，请换一个")
	}

	_, err = this.RPC().ServerRPC().UpdateServerDNSName(this.AdminContext(), &pb.UpdateServerDNSNameRequest{
		ServerId: params.ServerId,
		DnsName:  dnsName,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
