package certs

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"strings"
	"time"
)

// SelectPopupAction 选择证书
type SelectPopupAction struct {
	actionutils.ParentAction
}

func (this *SelectPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *SelectPopupAction) RunGet(params struct {
	ViewSize        string
	SelectedCertIds string
	Keyword         string
}) {
	// TODO 列出常用和最新的证书供用户选择

	this.Data["keyword"] = params.Keyword

	// 已经选择的证书
	selectedCertIds := []string{}
	if len(params.SelectedCertIds) > 0 {
		selectedCertIds = strings.Split(params.SelectedCertIds, ",")
	}

	if len(params.ViewSize) == 0 {
		params.ViewSize = "normal"
	}
	this.Data["viewSize"] = params.ViewSize

	countResp, err := this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
		Keyword: params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	page := this.NewPage(countResp.Count)
	this.Data["page"] = page.AsHTML()

	listResp, err := this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{
		Keyword: params.Keyword,
		Offset:  page.Offset,
		Size:    page.Size,
	})

	certConfigs := []*sslconfigs.SSLCertConfig{}
	err = json.Unmarshal(listResp.SslCertsJSON, &certConfigs)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["certs"] = certConfigs

	certMaps := []maps.Map{}
	nowTime := time.Now().Unix()
	for _, certConfig := range certConfigs {
		countServersResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithSSLCertId(this.AdminContext(), &pb.CountAllEnabledServersWithSSLCertIdRequest{SslCertId: certConfig.Id})
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
			"isSelected":   lists.ContainsString(selectedCertIds, numberutils.FormatInt64(certConfig.Id)),
		})
	}
	this.Data["certInfos"] = certMaps

	this.Show()
}
