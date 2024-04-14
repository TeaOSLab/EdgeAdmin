package http

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"regexp"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("http")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// 只有HTTP服务才支持
	if this.FilterHTTPFamily() {
		return
	}

	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}
	var httpConfig = &serverconfigs.HTTPProtocolConfig{}
	if len(server.HttpJSON) > 0 {
		err := json.Unmarshal(server.HttpJSON, httpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	} else {
		httpConfig.IsOn = true
	}
	_ = httpConfig.Init()
	var httpPorts = httpConfig.AllPorts()

	// 检查http和https端口冲突
	var conflictingPorts = []int{}
	if len(server.HttpsJSON) > 0 {
		var httpsConfig = &serverconfigs.HTTPSProtocolConfig{}
		err := json.Unmarshal(server.HttpsJSON, httpsConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		_ = httpsConfig.Init(context.TODO())
		for _, port := range httpsConfig.AllPorts() {
			if lists.ContainsInt(httpPorts, port) {
				conflictingPorts = append(conflictingPorts, port)
			}
		}
	}
	this.Data["conflictingPorts"] = conflictingPorts

	this.Data["serverType"] = server.Type
	this.Data["httpConfig"] = maps.Map{
		"isOn":      httpConfig.IsOn,
		"addresses": httpConfig.Listen,
	}

	// 跳转相关设置
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["webId"] = webConfig.Id
	this.Data["redirectToHTTPSConfig"] = webConfig.RedirectToHttps

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId  int64
	IsOn      bool
	Addresses string

	WebId               int64
	RedirectToHTTPSJSON []byte

	Must *actions.Must
}) {
	// 记录日志
	defer this.CreateLogInfo(codes.ServerHTTP_LogUpdateHTTPSettings, params.ServerId)

	var addresses = []*serverconfigs.NetworkAddressConfig{}
	err := json.Unmarshal([]byte(params.Addresses), &addresses)
	if err != nil {
		this.Fail("端口地址解析失败：" + err.Error())
	}

	// 如果启用HTTP时没有填写端口，则默认为80
	if params.IsOn && len(addresses) == 0 {
		addresses = []*serverconfigs.NetworkAddressConfig{
			{
				Protocol:  serverconfigs.ProtocolHTTP,
				PortRange: "80",
			},
		}
	}

	// 检查端口地址是否正确
	for _, addr := range addresses {
		err = addr.Init()
		if err != nil {
			this.Fail("绑定端口校验失败：" + err.Error())
		}

		if regexp.MustCompile(`^\d+$`).MatchString(addr.PortRange) {
			port := types.Int(addr.PortRange)
			if port > 65535 {
				this.Fail("绑定的端口地址不能大于65535")
			}
			if port == 443 {
				this.Fail("端口443通常是HTTPS的端口，不能用在HTTP上")
			}
		}
	}

	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}
	var httpConfig = &serverconfigs.HTTPProtocolConfig{}
	if len(server.HttpJSON) > 0 {
		err = json.Unmarshal(server.HttpJSON, httpConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	httpConfig.IsOn = params.IsOn
	httpConfig.Listen = addresses
	configData, err := json.Marshal(httpConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().ServerRPC().UpdateServerHTTP(this.AdminContext(), &pb.UpdateServerHTTPRequest{
		ServerId: params.ServerId,
		HttpJSON: configData,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 设置跳转到HTTPS
	// TODO 校验设置
	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebRedirectToHTTPS(this.AdminContext(), &pb.UpdateHTTPWebRedirectToHTTPSRequest{
		HttpWebId:           params.WebId,
		RedirectToHTTPSJSON: params.RedirectToHTTPSJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
