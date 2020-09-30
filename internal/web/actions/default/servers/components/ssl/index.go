package ssl

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.FirstMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	Type string
}) {
	this.Data["type"] = params.Type

	countAll := int64(0)
	countCA := int64(0)
	countAvailable := int64(0)
	countExpired := int64(0)
	count7Days := int64(0)
	count30Days := int64(0)

	// 计算数量
	{
		// all
		resp, err := this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countAll = resp.Count

		// CA
		resp, err = this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			IsCA: true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countCA = resp.Count

		// available
		resp, err = this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			IsAvailable: true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countAvailable = resp.Count

		// expired
		resp, err = this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			IsExpired: true,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countExpired = resp.Count

		// expire in 7 days
		resp, err = this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			ExpiringDays: 7,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		count7Days = resp.Count

		// expire in 30 days
		resp, err = this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			ExpiringDays: 30,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		count30Days = resp.Count
	}

	this.Data["countAll"] = countAll
	this.Data["countCA"] = countCA
	this.Data["countAvailable"] = countAvailable
	this.Data["countExpired"] = countExpired
	this.Data["count7Days"] = count7Days
	this.Data["count30Days"] = count30Days

	// 分页
	var page *actionutils.Page
	var listResp *pb.ListSSLCertsResponse
	var err error
	switch params.Type {
	case "":
		page = this.NewPage(countAll)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{Offset: page.Offset, Size: page.Size})
	case "ca":
		page = this.NewPage(countCA)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{IsCA: true, Offset: page.Offset, Size: page.Size})
	case "available":
		page = this.NewPage(countAvailable)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{IsAvailable: true, Offset: page.Offset, Size: page.Size})
	case "expired":
		page = this.NewPage(countExpired)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{IsExpired: true, Offset: page.Offset, Size: page.Size})
	case "7days":
		page = this.NewPage(count7Days)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{ExpiringDays: 7, Offset: page.Offset, Size: page.Size})
	case "30days":
		page = this.NewPage(count30Days)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{ExpiringDays: 30, Offset: page.Offset, Size: page.Size})
	default:
		page = this.NewPage(countAll)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{})
	}
	if err != nil {
		this.ErrorPage(err)
		return
	}

	certConfigs := []*sslconfigs.SSLCertConfig{}
	err = json.Unmarshal(listResp.CertsJSON, &certConfigs)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["certs"] = certConfigs

	certMaps := []maps.Map{}
	nowTime := time.Now().Unix()
	for _, certConfig := range certConfigs {
		countServersResp, err := this.RPC().ServerRPC().CountServersWithSSLCertId(this.AdminContext(), &pb.CountServersWithSSLCertIdRequest{CertId: certConfig.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		certMaps = append(certMaps, maps.Map{
			"beginDay":     timeutil.FormatTime("Y-m-d", certConfig.TimeBeginAt),
			"endDay":       timeutil.FormatTime("Y-m-d", certConfig.TimeEndAt),
			"isExpired":    nowTime > certConfig.TimeEndAt,
			"isAvailable":  nowTime <= certConfig.TimeEndAt,
			"countServers": countServersResp.Count,
		})
	}
	this.Data["certInfos"] = certMaps

	this.Data["page"] = page.AsHTML()

	this.Show()
}
