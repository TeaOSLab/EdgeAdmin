// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package origins

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"net"
	"net/url"
	"strings"
	"time"
)

type DetectHTTPSAction struct {
	actionutils.ParentAction
}

func (this *DetectHTTPSAction) RunPost(params struct {
	Addr string
}) {
	this.Data["isOk"] = false

	// parse from url
	if strings.HasPrefix(params.Addr, "http://") || strings.HasPrefix(params.Addr, "https://") {
		u, err := url.Parse(params.Addr)
		if err == nil {
			params.Addr = u.Host
		}
	}

	this.Data["addr"] = params.Addr

	if len(params.Addr) == 0 {
		this.Success()
		return
	}

	var realHost = params.Addr
	host, port, err := net.SplitHostPort(params.Addr)
	if err == nil {
		if port != "80" {
			this.Success()
			return
		}
		realHost = host
	}

	conn, err := net.DialTimeout("tcp", configutils.QuoteIP(realHost)+":443", 3*time.Second)
	if err != nil {
		this.Success()
		return
	}
	_ = conn.Close()

	this.Data["isOk"] = true

	this.Success()
}
