package server

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	adminserverutils "github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/settings/server/admin-server-utils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/iwind/TeaGo/actions"
	"net"
)

type UpdateHTTPPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateHTTPPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateHTTPPopupAction) RunGet(params struct{}) {
	serverConfig, err := adminserverutils.LoadServerConfig()
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["serverConfig"] = serverConfig

	this.Show()
}

func (this *UpdateHTTPPopupAction) RunPost(params struct {
	IsOn    bool
	Listens []string

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.AdminServer_LogUpdateServerHTTPSettings)

	if len(params.Listens) == 0 {
		this.Fail("请输入绑定地址")
	}

	serverConfig, err := adminserverutils.LoadServerConfig()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	serverConfig.Http.On = params.IsOn

	listen := []string{}
	for _, addr := range params.Listens {
		addr = utils.FormatAddress(addr)
		if len(addr) == 0 {
			continue
		}
		if _, _, err := net.SplitHostPort(addr); err != nil {
			addr += ":80"
		}
		listen = append(listen, addr)
	}
	serverConfig.Http.Listen = listen

	err = adminserverutils.WriteServerConfig(serverConfig)
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
