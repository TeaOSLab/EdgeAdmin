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

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.FirstMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	UserId   int64
	Type     string // [empty] | ca | 7days | ...
	Keyword  string
	UserType string
}) {
	this.Data["type"] = params.Type
	this.Data["keyword"] = params.Keyword

	if params.UserId > 0 {
		params.UserType = "user"
	}
	this.Data["userType"] = params.UserType

	// 当前用户
	this.Data["searchingUserId"] = params.UserId
	var userMap = maps.Map{
		"id":       0,
		"username": "",
		"fullname": "",
	}
	if params.UserId > 0 {
		userResp, err := this.RPC().UserRPC().FindEnabledUser(this.AdminContext(), &pb.FindEnabledUserRequest{UserId: params.UserId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var user = userResp.User
		if user != nil {
			userMap = maps.Map{
				"id":       user.Id,
				"username": user.Username,
				"fullname": user.Fullname,
			}
		}
	}
	this.Data["user"] = userMap

	var countAll int64
	var countCA int64
	var countAvailable int64
	var countExpired int64
	var count7Days int64
	var count30Days int64

	var userOnly = params.UserType == "user" || params.UserId > 0

	// 计算数量
	{
		// all
		resp, err := this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			UserId:   params.UserId,
			Keyword:  params.Keyword,
			UserOnly: userOnly,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countAll = resp.Count

		// CA
		resp, err = this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			UserId:   params.UserId,
			IsCA:     true,
			Keyword:  params.Keyword,
			UserOnly: userOnly,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countCA = resp.Count

		// available
		resp, err = this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			UserId:      params.UserId,
			IsAvailable: true,
			Keyword:     params.Keyword,
			UserOnly:    userOnly,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countAvailable = resp.Count

		// expired
		resp, err = this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			UserId:    params.UserId,
			IsExpired: true,
			Keyword:   params.Keyword,
			UserOnly:  userOnly,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countExpired = resp.Count

		// expire in 7 days
		resp, err = this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			UserId:       params.UserId,
			ExpiringDays: 7,
			Keyword:      params.Keyword,
			UserOnly:     userOnly,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		count7Days = resp.Count

		// expire in 30 days
		resp, err = this.RPC().SSLCertRPC().CountSSLCerts(this.AdminContext(), &pb.CountSSLCertRequest{
			UserId:       params.UserId,
			ExpiringDays: 30,
			Keyword:      params.Keyword,
			UserOnly:     userOnly,
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
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{
			UserId:   params.UserId,
			Offset:   page.Offset,
			Size:     page.Size,
			Keyword:  params.Keyword,
			UserOnly: userOnly,
		})
	case "ca":
		page = this.NewPage(countCA)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{
			UserId:   params.UserId,
			IsCA:     true,
			Offset:   page.Offset,
			Size:     page.Size,
			Keyword:  params.Keyword,
			UserOnly: userOnly,
		})
	case "available":
		page = this.NewPage(countAvailable)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{
			UserId:      params.UserId,
			IsAvailable: true,
			Offset:      page.Offset,
			Size:        page.Size,
			Keyword:     params.Keyword,
			UserOnly:    userOnly,
		})
	case "expired":
		page = this.NewPage(countExpired)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{
			UserId:    params.UserId,
			IsExpired: true,
			Offset:    page.Offset,
			Size:      page.Size,
			Keyword:   params.Keyword,
			UserOnly:  userOnly,
		})
	case "7days":
		page = this.NewPage(count7Days)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{
			UserId:       params.UserId,
			ExpiringDays: 7,
			Offset:       page.Offset,
			Size:         page.Size,
			Keyword:      params.Keyword,
		})
	case "30days":
		page = this.NewPage(count30Days)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{
			UserId:       params.UserId,
			ExpiringDays: 30,
			Offset:       page.Offset,
			Size:         page.Size,
			Keyword:      params.Keyword,
			UserOnly:     userOnly,
		})
	default:
		page = this.NewPage(countAll)
		listResp, err = this.RPC().SSLCertRPC().ListSSLCerts(this.AdminContext(), &pb.ListSSLCertsRequest{
			UserId:   params.UserId,
			Keyword:  params.Keyword,
			UserOnly: userOnly,
			Offset:   page.Offset,
			Size:     page.Size,
		})
	}
	if err != nil {
		this.ErrorPage(err)
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
		// count servers
		countServersResp, err := this.RPC().ServerRPC().CountAllEnabledServersWithSSLCertId(this.AdminContext(), &pb.CountAllEnabledServersWithSSLCertIdRequest{
			SslCertId: certConfig.Id,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// user
		userResp, err := this.RPC().SSLCertRPC().FindSSLCertUser(this.AdminContext(), &pb.FindSSLCertUserRequest{SslCertId: certConfig.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var certUserMap = maps.Map{
			"id": 0,
		}
		if userResp.User != nil {
			certUserMap = maps.Map{
				"id":       userResp.User.Id,
				"username": userResp.User.Username,
				"fullname": userResp.User.Fullname,
			}
		}

		certMaps = append(certMaps, maps.Map{
			"isOn":         certConfig.IsOn,
			"beginDay":     timeutil.FormatTime("Y-m-d", certConfig.TimeBeginAt),
			"endDay":       timeutil.FormatTime("Y-m-d", certConfig.TimeEndAt),
			"isExpired":    nowTime > certConfig.TimeEndAt,
			"isAvailable":  nowTime <= certConfig.TimeEndAt,
			"countServers": countServersResp.Count,
			"user":         certUserMap,
		})
	}
	this.Data["certInfos"] = certMaps

	this.Data["page"] = page.AsHTML()

	this.Show()
}
