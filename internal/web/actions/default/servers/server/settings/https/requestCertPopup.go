package https

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"strings"
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
	} else {
		this.Fail("请选择或者填写用户")
	}

	createTaskResp, err := this.RPC().ACMETaskRPC().CreateACMETask(this.AdminContext(), &pb.CreateACMETaskRequest{
		AcmeUserId:    acmeUserId,
		DnsProviderId: 0,
		DnsDomain:     "",
		Domains:       params.ServerNames,
		AutoRenew:     true,
		AuthType:      "http",
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
