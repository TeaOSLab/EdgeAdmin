package certs

import (
	"encoding/json"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
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
	ServerId         int64  // 搜索的服务
	UserId           int64  // 搜索的用户名
	SearchingDomains string // 搜索的域名
	SearchingType    string // 搜索类型：match|all

	ViewSize        string
	SelectedCertIds string
	Keyword         string
}) {
	this.Data["searchingServerId"] = params.ServerId

	// 服务相关
	if params.ServerId > 0 {
		serverResp, err := this.RPC().ServerRPC().FindEnabledUserServerBasic(this.AdminContext(), &pb.FindEnabledUserServerBasicRequest{ServerId: params.ServerId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var server = serverResp.Server
		if server != nil {
			if server.UserId > 0 {
				params.UserId = server.UserId
			}

			// 读取所有ServerNames
			serverNamesResp, err := this.RPC().ServerRPC().FindServerNames(this.AdminContext(), &pb.FindServerNamesRequest{ServerId: params.ServerId})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if len(serverNamesResp.ServerNamesJSON) > 0 {
				var serverNames = []*serverconfigs.ServerNameConfig{}
				err = json.Unmarshal(serverNamesResp.ServerNamesJSON, &serverNames)
				if err != nil {
					this.ErrorPage(err)
					return
				}
				params.SearchingDomains = strings.Join(serverconfigs.PlainServerNames(serverNames), ",")
			}
		}
	}

	// 用户相关
	this.Data["userId"] = params.UserId // 可变
	this.Data["searchingUserId"] = params.UserId

	// 域名搜索相关
	var url = this.Request.URL.Path
	var query = this.Request.URL.Query()
	query.Del("searchingType")
	this.Data["baseURL"] = url + "?" + query.Encode()

	var searchingDomains = []string{}
	if len(params.SearchingDomains) > 0 {
		searchingDomains = strings.Split(params.SearchingDomains, ",")
	}
	const maxDomains = 2_000 // 限制搜索的域名数量
	if len(searchingDomains) > maxDomains {
		searchingDomains = searchingDomains[:maxDomains]
	}
	this.Data["allSearchingDomains"] = params.SearchingDomains
	this.Data["searchingDomains"] = searchingDomains

	this.Data["keyword"] = params.Keyword
	this.Data["selectedCertIds"] = params.SelectedCertIds

	var searchingType = params.SearchingType
	if len(searchingType) == 0 {
		if len(params.SearchingDomains) == 0 {
			searchingType = "all"
		} else {
			searchingType = "match"
		}
	}
	if searchingType != "all" && searchingType != "match" {
		this.ErrorPage(errors.New("invalid searching type '" + searchingType + "'"))
		return
	}
	this.Data["searchingType"] = searchingType

	// 已经选择的证书
	var selectedCertIds = []string{}
	if len(params.SelectedCertIds) > 0 {
		selectedCertIds = strings.Split(params.SelectedCertIds, ",")
	}

	if len(params.ViewSize) == 0 {
		params.ViewSize = "normal"
	}
	this.Data["viewSize"] = params.ViewSize

	// 全部证书数量
	countAllResp, err := this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
		UserId:  params.UserId,
		Keyword: params.Keyword,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var totalAll = countAllResp.Count
	this.Data["totalAll"] = totalAll

	// 已匹配证书数量
	var totalMatch int64 = 0
	if len(searchingDomains) > 0 {
		countMatchResp, err := this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			UserId:  params.UserId,
			Keyword: params.Keyword,
			Domains: searchingDomains,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		totalMatch = countMatchResp.Count
	}
	this.Data["totalMatch"] = totalMatch

	var totalCerts int64
	if searchingType == "all" {
		totalCerts = totalAll
	} else if searchingType == "match" {
		totalCerts = totalMatch
	}

	var page = this.NewPage(totalCerts)
	this.Data["page"] = page.AsHTML()

	var listResp *pb.ListSSLCertsResponse
	if searchingType == "all" {
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{
			UserId:  params.UserId,
			Keyword: params.Keyword,
			Offset:  page.Offset,
			Size:    page.Size,
		})
	} else if searchingType == "match" {
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{
			UserId:  params.UserId,
			Keyword: params.Keyword,
			Domains: searchingDomains,
			Offset:  page.Offset,
			Size:    page.Size,
		})
	}

	if err != nil {
		this.ErrorPage(err)
		return
	}

	if listResp == nil {
		this.ErrorPage(errors.New("'listResp' should not be nil"))
		return
	}

	var certConfigs = []*sslconfigs.SSLCertConfig{}
	err = json.Unmarshal(listResp.SslCertsJSON, &certConfigs)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["certs"] = certConfigs

	var certMaps = []maps.Map{}
	var nowTime = time.Now().Unix()
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
