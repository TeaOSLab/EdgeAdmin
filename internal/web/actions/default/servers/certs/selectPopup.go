package certs

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"time"
)

// 选择证书
type SelectPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectPopupAction) RunGet(params struct {
	ViewSize string
}) {
	// TODO 支持关键词搜索
	// TODO 列出常用的证书供用户选择

	if len(params.ViewSize) == 0 {
		params.ViewSize = "normal"
	}
	this.Data["viewSize"] = params.ViewSize

	countResp, err := this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	page := this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	listResp, err := this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})

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
		countServersResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithSSLCertId(this.AdminContext(), &pb.CountAllEnabledServersWithSSLCertIdRequest{CertId: certConfig.Id})
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

	this.Show()
}
