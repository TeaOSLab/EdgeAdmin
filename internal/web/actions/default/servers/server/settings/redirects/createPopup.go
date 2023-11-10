package redirects

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"net/url"
	"regexp"
	"strings"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
}) {
	this.Data["statusList"] = serverconfigs.AllHTTPRedirectStatusList()

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Type string

	// URL
	Mode           string
	BeforeURL      string
	AfterURL       string
	MatchPrefix    bool
	MatchRegexp    bool
	KeepRequestURI bool
	KeepArgs       bool

	// 域名
	DomainsAll              bool
	DomainsBeforeJSON       []byte
	DomainBeforeIgnorePorts bool
	DomainAfter             string
	DomainAfterScheme       string

	// 端口
	PortsAll        bool
	PortsBefore     []string
	PortAfter       int
	PortAfterScheme string

	Status int

	ExceptDomainsJSON []byte
	OnlyDomainsJSON   []byte

	CondsJSON []byte
	IsOn      bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var config = &serverconfigs.HTTPHostRedirectConfig{}
	config.Type = params.Type
	config.Status = params.Status
	config.IsOn = params.IsOn

	switch params.Type {
	case serverconfigs.HTTPHostRedirectTypeURL:
		params.Must.
			Field("beforeURL", params.BeforeURL).
			Require("请填写跳转前的URL")

		// 校验格式
		if params.MatchRegexp {
			_, err := regexp.Compile(params.BeforeURL)
			if err != nil {
				this.Fail("跳转前URL正则表达式错误：" + err.Error())
			}
		} else {
			u, err := url.Parse(params.BeforeURL)
			if err != nil {
				this.FailField("beforeURL", "请输入正确的跳转前URL")
			}
			if (u.Scheme != "http" && u.Scheme != "https") ||
				len(u.Host) == 0 {
				this.FailField("beforeURL", "请输入正确的跳转前URL")
			}

		}

		params.Must.
			Field("afterURL", params.AfterURL).
			Require("请填写跳转后URL")

		// 校验格式
		if params.MatchRegexp {
			// 正则表达式情况下不做校验
		} else {
			u, err := url.Parse(params.AfterURL)
			if err != nil {
				this.FailField("afterURL", "请输入正确的跳转后URL")
			}
			if (u.Scheme != "http" && u.Scheme != "https") ||
				len(u.Host) == 0 {
				this.FailField("afterURL", "请输入正确的跳转后URL")
			}
		}

		config.Mode = params.Mode
		config.BeforeURL = params.BeforeURL
		config.AfterURL = params.AfterURL
		config.MatchPrefix = params.MatchPrefix
		config.MatchRegexp = params.MatchRegexp
		config.KeepRequestURI = params.KeepRequestURI
		config.KeepArgs = params.KeepArgs
	case serverconfigs.HTTPHostRedirectTypeDomain:
		config.DomainsAll = params.DomainsAll
		var domainsBefore = []string{}
		if len(params.DomainsBeforeJSON) > 0 {
			err := json.Unmarshal(params.DomainsBeforeJSON, &domainsBefore)
			if err != nil {
				this.Fail("错误的域名格式：" + err.Error())
				return
			}
		}
		config.DomainsBefore = domainsBefore
		if !params.DomainsAll {
			if len(domainsBefore) == 0 {
				this.Fail("请输入跳转前域名")
				return
			}
		}
		config.DomainBeforeIgnorePorts = params.DomainBeforeIgnorePorts
		if len(params.DomainAfter) == 0 {
			this.FailField("domainAfter", "请输入跳转后域名")
			return
		}

		// 检查用户输入的是否为域名
		if !domainutils.ValidateDomainFormat(params.DomainAfter) {
			// 是否为URL
			u, err := url.Parse(params.DomainAfter)
			if err == nil {
				if len(u.Host) == 0 {
					this.FailField("domainAfter", "跳转后域名输入不正确")
					return
				}
				params.DomainAfter = u.Host
			} else {
				this.FailField("domainAfter", "跳转后域名输入不正确")
				return
			}
		}

		config.DomainAfter = params.DomainAfter
		config.DomainAfterScheme = params.DomainAfterScheme
	case serverconfigs.HTTPHostRedirectTypePort:
		config.PortsAll = params.PortsAll

		config.PortsBefore = params.PortsBefore
		var portReg = regexp.MustCompile(`^\d+$`)
		var portRangeReg = regexp.MustCompile(`^\d+-\d+$`)
		if !config.PortsAll {
			for _, port := range params.PortsBefore {
				port = strings.ReplaceAll(port, " ", "")
				if !portReg.MatchString(port) && !portRangeReg.MatchString(port) {
					this.Fail("端口号" + port + "填写错误（请输入单个端口号或一个端口范围）")
					return
				}
			}
			if len(params.PortsBefore) == 0 {
				this.Fail("请输入跳转前端口")
				return
			}
		}

		if params.PortAfter <= 0 {
			this.FailField("portAfter", "请输入跳转后端口")
			return
		}
		config.PortAfter = params.PortAfter
		config.PortAfterScheme = params.PortAfterScheme
	}

	params.Must.
		Field("status", params.Status).
		Gte(0, "请选择正确的跳转状态码")

	// 域名
	if len(params.ExceptDomainsJSON) > 0 {
		var exceptDomains = []string{}
		err := json.Unmarshal(params.ExceptDomainsJSON, &exceptDomains)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		config.ExceptDomains = exceptDomains
	}

	if len(params.OnlyDomainsJSON) > 0 {
		var onlyDomains = []string{}
		err := json.Unmarshal(params.OnlyDomainsJSON, &onlyDomains)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		config.OnlyDomains = onlyDomains
	}

	// 校验匹配条件
	var conds *shared.HTTPRequestCondsConfig
	if len(params.CondsJSON) > 0 {
		conds = &shared.HTTPRequestCondsConfig{}
		err := json.Unmarshal(params.CondsJSON, conds)
		if err != nil {
			this.Fail("匹配条件校验失败：" + err.Error())
		}

		err = conds.Init()
		if err != nil {
			this.Fail("匹配条件校验失败：" + err.Error())
		}
	}
	config.Conds = conds

	// 校验配置
	err := config.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
		return
	}

	this.Data["redirect"] = config

	this.Success()
}
