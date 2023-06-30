package servers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net/url"
	"regexp"
	"strings"
)

type AddOriginPopupAction struct {
	actionutils.ParentAction
}

func (this *AddOriginPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *AddOriginPopupAction) RunGet(params struct {
	ServerType string
}) {
	this.Data["serverType"] = params.ServerType

	this.getOSSHook()

	this.Show()
}

func (this *AddOriginPopupAction) RunPost(params struct {
	Protocol string
	Addr     string

	DomainsJSON  []byte
	Host         string
	FollowPort   bool
	Http2Enabled bool

	Must *actions.Must
}) {
	ossConfig, goNext, err := this.postOSSHook(params.Protocol)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if !goNext {
		return
	}

	// 初始化
	var pbAddr = &pb.NetworkAddress{
		Protocol: params.Protocol,
	}
	var addrConfig = &serverconfigs.NetworkAddressConfig{
		Protocol: serverconfigs.Protocol(params.Protocol),
	}
	var ossJSON []byte

	if ossConfig != nil { // OSS
		ossJSON, err = json.Marshal(ossConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		err = ossConfig.Init()
		if err != nil {
			this.Fail("校验OSS配置时出错：" + err.Error())
			return
		}
	} else { // 普通源站
		params.Must.
			Field("addr", params.Addr).
			Require("请输入源站地址")

		var addr = params.Addr

		// 是否是完整的地址
		if (params.Protocol == "http" || params.Protocol == "https") && regexp.MustCompile(`^(http|https)://`).MatchString(addr) {
			u, err := url.Parse(addr)
			if err == nil {
				addr = u.Host
			}
		}

		addr = regexp.MustCompile(`\s+`).ReplaceAllString(addr, "")
		var portIndex = strings.LastIndex(addr, ":")
		if portIndex < 0 {
			if params.Protocol == "http" {
				addr += ":80"
			} else if params.Protocol == "https" {
				addr += ":443"
			} else {
				this.Fail("地址中需要带有端口")
			}
			portIndex = strings.LastIndex(addr, ":")
		}
		var host = addr[:portIndex]
		var port = addr[portIndex+1:]

		// 检查端口号
		if port == "0" {
			this.Fail("端口号不能为0")
		}
		if !configutils.HasVariables(port) {
			// 必须是整数
			if !regexp.MustCompile(`^\d+$`).MatchString(port) {
				this.Fail("端口号只能为整数")
			}
			var portInt = types.Int(port)
			if portInt == 0 {
				this.Fail("端口号不能为0")
			}
			if portInt > 65535 {
				this.Fail("端口号不能大于65535")
			}
		}

		pbAddr = &pb.NetworkAddress{
			Protocol:  params.Protocol,
			Host:      host,
			PortRange: port,
		}

		addrConfig = &serverconfigs.NetworkAddressConfig{
			Protocol:  serverconfigs.Protocol(params.Protocol),
			Host:      host,
			PortRange: port,
		}
	}

	// 专属域名
	var domains = []string{}
	if len(params.DomainsJSON) > 0 {
		err := json.Unmarshal(params.DomainsJSON, &domains)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		// 去除可能误加的斜杠
		for index, domain := range domains {
			domains[index] = strings.TrimSuffix(domain, "/")
		}
	}

	resp, err := this.RPC().OriginRPC().CreateOrigin(this.AdminContext(), &pb.CreateOriginRequest{
		Name:         "",
		Addr:         pbAddr,
		OssJSON:      ossJSON,
		Description:  "",
		Weight:       10,
		IsOn:         true,
		Domains:      domains,
		Host:         params.Host,
		FollowPort:   params.FollowPort,
		Http2Enabled: params.Http2Enabled,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var origin = &serverconfigs.OriginConfig{
		Id:   resp.OriginId,
		IsOn: true,
		Addr: addrConfig,
		OSS:  ossConfig,
	}

	this.Data["origin"] = maps.Map{
		"id":          resp.OriginId,
		"isOn":        true,
		"addr":        addrConfig,
		"addrSummary": origin.AddrSummary(),
	}

	// 创建日志
	defer this.CreateLogInfo(codes.ServerOrigin_LogCreateOrigin, resp.OriginId)

	this.Success()
}
