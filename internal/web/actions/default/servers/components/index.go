package components

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

const (
	SettingCodeServerGlobalConfig = "serverGlobalConfig"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "component", "index")
	this.SecondMenu("global")
}

func (this *IndexAction) RunGet(params struct{}) {
	valueJSONResp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: SettingCodeServerGlobalConfig})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	valueJSON := valueJSONResp.ValueJSON
	globalConfig := &serverconfigs.GlobalConfig{}

	// 默认值
	globalConfig.HTTPAll.DomainAuditingIsOn = true

	if len(valueJSON) > 0 {
		err = json.Unmarshal(valueJSON, globalConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["globalConfig"] = globalConfig

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	GlobalConfigJSON []byte
	Must             *actions.Must

	// 不匹配域名相关
	AllowMismatchDomains                []string
	DomainMismatchAction                string
	DomainMismatchActionPageStatusCode  int
	DomainMismatchActionPageContentHTML string

	// TCP端口设置
	TcpAllPortRangeMin int
	TcpAllPortRangeMax int
	TcpAllDenyPorts    []int

	DefaultDomain string
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "保存代理服务全局配置")

	if len(params.GlobalConfigJSON) == 0 {
		this.Fail("错误的配置信息，请刷新当前页面后重试")
	}

	globalConfig := &serverconfigs.GlobalConfig{}
	err := json.Unmarshal(params.GlobalConfigJSON, globalConfig)
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
	}

	// 允许不匹配的域名
	allowMismatchDomains := []string{}
	for _, domain := range params.AllowMismatchDomains {
		if len(domain) > 0 {
			allowMismatchDomains = append(allowMismatchDomains, domain)
		}
	}
	globalConfig.HTTPAll.AllowMismatchDomains = allowMismatchDomains

	// 不匹配域名的动作
	switch params.DomainMismatchAction {
	case "close":
		globalConfig.HTTPAll.DomainMismatchAction = &serverconfigs.DomainMismatchAction{
			Code:    "close",
			Options: nil,
		}
	case "page":
		if params.DomainMismatchActionPageStatusCode <= 0 {
			params.DomainMismatchActionPageStatusCode = 404
		}
		globalConfig.HTTPAll.DomainMismatchAction = &serverconfigs.DomainMismatchAction{
			Code: "page",
			Options: maps.Map{
				"statusCode":  params.DomainMismatchActionPageStatusCode,
				"contentHTML": params.DomainMismatchActionPageContentHTML,
			},
		}
	}

	// TCP端口范围
	if params.TcpAllPortRangeMin < 1024 {
		params.TcpAllPortRangeMin = 1024
	}
	if params.TcpAllPortRangeMax > 65534 {
		params.TcpAllPortRangeMax = 65534
	} else if params.TcpAllPortRangeMax < 1024 {
		params.TcpAllPortRangeMax = 1024
	}
	if params.TcpAllPortRangeMin > params.TcpAllPortRangeMax {
		params.TcpAllPortRangeMin, params.TcpAllPortRangeMax = params.TcpAllPortRangeMax, params.TcpAllPortRangeMin
	}
	globalConfig.TCPAll.DenyPorts = params.TcpAllDenyPorts
	globalConfig.TCPAll.PortRangeMin = params.TcpAllPortRangeMin
	globalConfig.TCPAll.PortRangeMax = params.TcpAllPortRangeMax

	// 修改配置
	globalConfigJSON, err := json.Marshal(globalConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().SysSettingRPC().UpdateSysSetting(this.AdminContext(), &pb.UpdateSysSettingRequest{
		Code:      SettingCodeServerGlobalConfig,
		ValueJSON: globalConfigJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	//  通知更新
	_, err = this.RPC().ServerRPC().NotifyServersChange(this.AdminContext(), &pb.NotifyServersChangeRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
