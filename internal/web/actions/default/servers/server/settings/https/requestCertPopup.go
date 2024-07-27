package https

import (
	"encoding/json"
	"fmt"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"math/rand"
	"strings"
	"sync"

	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type RequestCertPopupAction struct {
	actionutils.ParentAction
}

func (this *RequestCertPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *RequestCertPopupAction) RunGet(params struct {
	ServerId           int64
	ExcludeServerNames string
}) {
	serverNamesResp, err := this.RPC().ServerRPC().FindServerNames(this.AdminContext(), &pb.FindServerNamesRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var serverNameConfigs = []*serverconfigs.ServerNameConfig{}
	err = json.Unmarshal(serverNamesResp.ServerNamesJSON, &serverNameConfigs)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var excludeServerNames = []string{}
	if len(params.ExcludeServerNames) > 0 {
		excludeServerNames = strings.Split(params.ExcludeServerNames, ",")
	}
	var serverNames = []string{}
	for _, c := range serverNameConfigs {
		if len(c.SubNames) == 0 {
			if domainutils.ValidateDomainFormat(c.Name) && !lists.ContainsString(excludeServerNames, c.Name) {
				serverNames = append(serverNames, c.Name)
			}
		} else {
			for _, subName := range c.SubNames {
				if domainutils.ValidateDomainFormat(subName) && !lists.ContainsString(excludeServerNames, subName) {
					serverNames = append(serverNames, subName)
				}
			}
		}
	}
	this.Data["serverNames"] = serverNames

	// 用户
	acmeUsersResp, err := this.RPC().ACMEUserRPC().FindAllACMEUsers(this.AdminContext(), &pb.FindAllACMEUsersRequest{
		AdminId: this.AdminId(),
		UserId:  0,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var userMaps = []maps.Map{}
	for _, user := range acmeUsersResp.AcmeUsers {
		description := user.Description
		if len(description) > 0 {
			description = "（" + description + "）"
		}

		userMaps = append(userMaps, maps.Map{
			"id":          user.Id,
			"description": description,
			"email":       user.Email,
		})
	}
	this.Data["users"] = userMaps

	this.Show()
}

func (this *RequestCertPopupAction) RunPost(params struct {
	ServerNames []string

	UserId    int64
	UserEmail string

	AsyncCreateCert bool
	ProviderCode    string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	// 检查域名
	if len(params.ServerNames) == 0 {
		this.Fail("必须包含至少一个或多个域名")
	}

	// 注册用户
	var acmeUserId int64
	if params.UserId > 0 {
		// TODO 检查当前管理员是否可以使用此用户
		acmeUserId = params.UserId
	} else if len(params.UserEmail) > 0 {
		params.Must.
			Field("userEmail", params.UserEmail).
			Email("Email格式错误")

		createUserResp, err := this.RPC().ACMEUserRPC().CreateACMEUser(this.AdminContext(), &pb.CreateACMEUserRequest{
			Email:       params.UserEmail,
			Description: "",
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		defer this.CreateLogInfo(codes.ACMEUser_LogCreateACMEUser, createUserResp.AcmeUserId)
		acmeUserId = createUserResp.AcmeUserId

		this.Data["acmeUser"] = maps.Map{
			"id":    acmeUserId,
			"email": params.UserEmail,
		}
	} else if !params.AsyncCreateCert {
		this.Fail("请选择或者填写用户")
	}

	if params.AsyncCreateCert && acmeUserId == 0 && params.ProviderCode != "" {
		acmeUsersResp, err := this.RPC().ACMEUserRPC().FindAllACMEUsers(this.AdminContext(), &pb.FindAllACMEUsersRequest{})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var userIds = []int64{}
		for _, user := range acmeUsersResp.AcmeUsers {
			if user.AcmeProviderCode == params.ProviderCode {
				userIds = append(userIds, user.Id)
			}
		}
		if len(userIds) < 5 {
			for i := 0; i < 5-len(userIds); i++ {
				createUserResp, err := this.RPC().ACMEUserRPC().CreateACMEUser(this.AdminContext(), &pb.CreateACMEUserRequest{
					Email:            fmt.Sprintf("%s-acme-%d-@mail.com", params.ProviderCode, rand.Intn(1000)),
					AcmeProviderCode: params.ProviderCode,
				})
				if err != nil {
					this.ErrorPage(err)
					return
				}
				userIds = append(userIds, createUserResp.AcmeUserId)
			}
		}

		acmeUserId = userIds[0]

		// 根据主域名拆分
		domainMap := make(map[string][]string)
		for _, domain := range params.ServerNames {
			parts := strings.Split(domain, ".")
			rootDomain := strings.Join(parts[len(parts)-2:], ".")
			domainMap[rootDomain] = append(domainMap[rootDomain], domain)
		}

		domainsGroup := make([][]string, 0, len(domainMap))
		for _, domains := range domainMap {
			domainsGroup = append(domainsGroup, domains)
		}

		domainsCount := len(domainsGroup)
		var queue = make(chan []string, domainsCount)
		for _, domains := range domainsGroup {
			queue <- domains
		}
		var concurrent = 10
		var wg = sync.WaitGroup{}
		wg.Add(domainsCount)
		sem := make(chan struct{}, concurrent)
		for i := 0; i < domainsCount; i++ {
			domains := <-queue
			sem <- struct{}{}
			go func(domains []string) {
				defer func() {
					<-sem
					wg.Done()
				}()
				createResp, err := this.RPC().ACMETaskRPC().CreateACMETask(this.AdminContext(), &pb.CreateACMETaskRequest{
					AcmeUserId:    acmeUserId,
					DnsProviderId: 0,
					DnsDomain:     "",
					Domains:       domains,
					AutoRenew:     true,
					AuthType:      "http",
					Async:         true,
				})
				if err == nil {
					this.CreateLogInfo("创建证书申请任务 %d", createResp.AcmeTaskId)
				}
			}(domains)
		}
		wg.Wait()
		this.Success()
	} else {
		createTaskResp, err := this.RPC().ACMETaskRPC().CreateACMETask(this.AdminContext(), &pb.CreateACMETaskRequest{
			AcmeUserId:    acmeUserId,
			DnsProviderId: 0,
			DnsDomain:     "",
			Domains:       params.ServerNames,
			AutoRenew:     true,
			AuthType:      "http",
			Async:         false,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		taskId := createTaskResp.AcmeTaskId

		defer this.CreateLogInfo(codes.ACMETask_LogRunACMETask, taskId)

		runResp, err := this.RPC().ACMETaskRPC().RunACMETask(this.AdminContext(), &pb.RunACMETaskRequest{AcmeTaskId: taskId})
		if err != nil {
			this.ErrorPage(err)
			return
		}

		if runResp.IsOk {
			certId := runResp.SslCertId

			configResp, err := this.RPC().SSLCertRPC().FindEnabledSSLCertConfig(this.AdminContext(), &pb.FindEnabledSSLCertConfigRequest{SslCertId: certId})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			certConfig := &sslconfigs.SSLCertConfig{}
			err = json.Unmarshal(configResp.SslCertJSON, certConfig)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			certConfig.CertData = nil // 去掉不必要的数据
			certConfig.KeyData = nil  // 去掉不必要的数据
			this.Data["cert"] = certConfig
			this.Data["certRef"] = &sslconfigs.SSLCertRef{
				IsOn:   true,
				CertId: certId,
			}

			this.Success()
		} else {
			this.Fail(runResp.Error)
		}
	}
}
