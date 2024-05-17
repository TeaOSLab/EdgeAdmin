// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package ocsp

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.SecondMenu("ocsp")
}

func (this *IndexAction) RunGet(params struct {
	Keyword string
}) {
	this.Data["keyword"] = params.Keyword

	countResp, err := this.RPC().SSLCertRPC().CountAllSSLCertsWithOCSPError(this.AdminContext(), &pb.CountAllSSLCertsWithOCSPErrorRequest{Keyword: params.Keyword})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var count = countResp.Count
	var page = this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	certsResp, err := this.RPC().SSLCertRPC().ListSSLCertsWithOCSPError(this.AdminContext(), &pb.ListSSLCertsWithOCSPErrorRequest{
		Keyword: params.Keyword,
		Offset:  page.Offset,
		Size:    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var certMaps = []maps.Map{}
	for _, cert := range certsResp.SslCerts {
		certMaps = append(certMaps, maps.Map{
			"id":            cert.Id,
			"isOn":          cert.IsOn,
			"dnsNames":      cert.DnsNames,
			"commonNames":   cert.CommonNames,
			"hasOCSP":       len(cert.Ocsp) > 0,
			"ocspIsUpdated": cert.OcspIsUpdated,
			"ocspError":     cert.OcspError,
			"isCA":          cert.IsCA,
			"isACME":        cert.IsACME,
			"name":          cert.Name,
			"isExpired":     cert.TimeEndAt < time.Now().Unix(),
			"beginDay":      timeutil.FormatTime("Y-m-d", cert.TimeBeginAt),
			"endDay":        timeutil.FormatTime("Y-m-d", cert.TimeEndAt),
		})
	}
	this.Data["certs"] = certMaps

	this.Show()
}
