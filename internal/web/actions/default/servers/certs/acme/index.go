package acme

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "task")
	this.SecondMenu("list")
}

func (this *IndexAction) RunGet(params struct {
	UserId   int64
	Type     string
	Keyword  string
	UserType string
}) {
	this.Data["type"] = params.Type
	this.Data["keyword"] = params.Keyword
	this.Data["userType"] = params.UserType

	var userOnly = params.UserId > 0 || params.UserType == "user"

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
	var countAvailable int64
	var countExpired int64
	var count7Days int64
	var count30Days int64

	// 计算数量
	{
		// all
		resp, err := this.RPC().ACMETaskRPC().CountAllEnabledACMETasks(this.AdminContext(), &pb.CountAllEnabledACMETasksRequest{
			UserId:   params.UserId,
			Keyword:  params.Keyword,
			UserOnly: userOnly,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		countAll = resp.Count

		// available
		resp, err = this.RPC().ACMETaskRPC().CountAllEnabledACMETasks(this.AdminContext(), &pb.CountAllEnabledACMETasksRequest{
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
		resp, err = this.RPC().ACMETaskRPC().CountAllEnabledACMETasks(this.AdminContext(), &pb.CountAllEnabledACMETasksRequest{
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
		resp, err = this.RPC().ACMETaskRPC().CountAllEnabledACMETasks(this.AdminContext(), &pb.CountAllEnabledACMETasksRequest{
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
		resp, err = this.RPC().ACMETaskRPC().CountAllEnabledACMETasks(this.AdminContext(), &pb.CountAllEnabledACMETasksRequest{
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
	this.Data["countAvailable"] = countAvailable
	this.Data["countExpired"] = countExpired
	this.Data["count7Days"] = count7Days
	this.Data["count30Days"] = count30Days

	// 分页
	var page *actionutils.Page
	var tasksResp *pb.ListEnabledACMETasksResponse
	var err error
	switch params.Type {
	case "":
		page = this.NewPage(countAll)
		tasksResp, err = this.RPC().ACMETaskRPC().ListEnabledACMETasks(this.AdminContext(), &pb.ListEnabledACMETasksRequest{
			UserId:   params.UserId,
			Offset:   page.Offset,
			Size:     page.Size,
			Keyword:  params.Keyword,
			UserOnly: userOnly,
		})
	case "available":
		page = this.NewPage(countAvailable)
		tasksResp, err = this.RPC().ACMETaskRPC().ListEnabledACMETasks(this.AdminContext(), &pb.ListEnabledACMETasksRequest{
			UserId:      params.UserId,
			IsAvailable: true,
			Offset:      page.Offset,
			Size:        page.Size,
			Keyword:     params.Keyword,
			UserOnly:    userOnly,
		})
	case "expired":
		page = this.NewPage(countExpired)
		tasksResp, err = this.RPC().ACMETaskRPC().ListEnabledACMETasks(this.AdminContext(), &pb.ListEnabledACMETasksRequest{
			UserId:    params.UserId,
			IsExpired: true,
			Offset:    page.Offset,
			Size:      page.Size,
			Keyword:   params.Keyword,
			UserOnly:  userOnly,
		})
	case "7days":
		page = this.NewPage(count7Days)
		tasksResp, err = this.RPC().ACMETaskRPC().ListEnabledACMETasks(this.AdminContext(), &pb.ListEnabledACMETasksRequest{
			UserId:       params.UserId,
			ExpiringDays: 7,
			Offset:       page.Offset,
			Size:         page.Size,
			Keyword:      params.Keyword,
			UserOnly:     userOnly,
		})
	case "30days":
		page = this.NewPage(count30Days)
		tasksResp, err = this.RPC().ACMETaskRPC().ListEnabledACMETasks(this.AdminContext(), &pb.ListEnabledACMETasksRequest{
			UserId:       params.UserId,
			ExpiringDays: 30,
			Offset:       page.Offset,
			Size:         page.Size,
			Keyword:      params.Keyword,
			UserOnly:     userOnly,
		})
	default:
		page = this.NewPage(countAll)
		tasksResp, err = this.RPC().ACMETaskRPC().ListEnabledACMETasks(this.AdminContext(), &pb.ListEnabledACMETasksRequest{
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

	this.Data["page"] = page.AsHTML()

	var taskMaps = []maps.Map{}
	for _, task := range tasksResp.AcmeTasks {
		if task.AcmeUser == nil {
			continue
		}

		// 服务商
		var providerMap maps.Map
		if task.AcmeUser.AcmeProvider != nil {
			providerMap = maps.Map{
				"name": task.AcmeUser.AcmeProvider.Name,
				"code": task.AcmeUser.AcmeProvider.Code,
			}
		}

		// 账号
		var accountMap maps.Map
		if task.AcmeUser.AcmeProviderAccount != nil {
			accountMap = maps.Map{
				"id":   task.AcmeUser.AcmeProviderAccount.Id,
				"name": task.AcmeUser.AcmeProviderAccount.Name,
			}
		}

		// DNS服务商
		dnsProviderMap := maps.Map{}
		if task.AuthType == "dns" && task.DnsProvider != nil {
			dnsProviderMap = maps.Map{
				"id":   task.DnsProvider.Id,
				"name": task.DnsProvider.Name,
			}
		}

		// 证书
		var certMap maps.Map = nil
		if task.SslCert != nil {
			certMap = maps.Map{
				"id":        task.SslCert.Id,
				"name":      task.SslCert.Name,
				"beginTime": timeutil.FormatTime("Y-m-d", task.SslCert.TimeBeginAt),
				"endTime":   timeutil.FormatTime("Y-m-d", task.SslCert.TimeEndAt),
			}
		}

		// 日志
		var logMap maps.Map = nil
		if task.LatestACMETaskLog != nil {
			logMap = maps.Map{
				"id":          task.LatestACMETaskLog.Id,
				"isOk":        task.LatestACMETaskLog.IsOk,
				"error":       task.LatestACMETaskLog.Error,
				"createdTime": timeutil.FormatTime("m-d", task.CreatedAt),
			}
		}

		// user
		userResp, err := this.RPC().ACMETaskRPC().FindACMETaskUser(this.AdminContext(), &pb.FindACMETaskUserRequest{AcmeTaskId: task.Id})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var taskUserMap = maps.Map{
			"id": 0,
		}
		if userResp.User != nil {
			taskUserMap = maps.Map{
				"id":       userResp.User.Id,
				"username": userResp.User.Username,
				"fullname": userResp.User.Fullname,
			}
		}

		taskMaps = append(taskMaps, maps.Map{
			"id":       task.Id,
			"authType": task.AuthType,
			"acmeUser": maps.Map{
				"id":       task.AcmeUser.Id,
				"email":    task.AcmeUser.Email,
				"provider": providerMap,
				"account":  accountMap,
			},
			"dnsProvider": dnsProviderMap,
			"dnsDomain":   task.DnsDomain,
			"domains":     task.Domains,
			"autoRenew":   task.AutoRenew,
			"cert":        certMap,
			"log":         logMap,
			"user":        taskUserMap,
		})
	}
	this.Data["tasks"] = taskMaps

	this.Show()
}
