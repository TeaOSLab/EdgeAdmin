// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package health

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"net"
	"strings"
)

type CheckDomainAction struct {
	actionutils.ParentAction
}

func (this *CheckDomainAction) RunPost(params struct {
	Host      string
	ClusterId int64
}) {
	this.Data["isOk"] = true // 默认为TRUE

	var host = params.Host
	if len(host) > 0 &&
		!strings.Contains(host, "{") /** 包含变量 **/ {
		h, _, err := net.SplitHostPort(host)
		if err == nil && len(h) > 0 {
			host = h
		}

		// 是否为IP
		if net.ParseIP(host) != nil {
			this.Success()
			return
		}

		host = strings.ToLower(host)
		resp, err := this.RPC().ServerRPC().CheckServerNameDuplicationInNodeCluster(this.AdminContext(), &pb.CheckServerNameDuplicationInNodeClusterRequest{
			NodeClusterId:   params.ClusterId,
			ServerNames:     []string{host},
			SupportWildcard: true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(resp.DuplicatedServerNames) == 0 {
			this.Data["isOk"] = false
		}
	}

	this.Success()
}
