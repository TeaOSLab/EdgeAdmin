package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"regexp"
	"strings"
)

type AddPortPopupAction struct {
	actionutils.ParentAction
}

func (this *AddPortPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *AddPortPopupAction) RunGet(params struct {
	ServerType string
	Protocol   string
	From       string
}) {
	this.Data["from"] = params.From

	protocols := serverconfigs.AllServerProtocolsForType(params.ServerType)
	if len(params.Protocol) > 0 {
		result := []maps.Map{}
		for _, p := range protocols {
			if p.GetString("code") == params.Protocol {
				result = append(result, p)
			}
		}
		protocols = result
	}
	this.Data["protocols"] = protocols

	this.Show()
}

func (this *AddPortPopupAction) RunPost(params struct {
	Protocol string
	Address  string

	Must *actions.Must
}) {
	// 校验地址
	addr := maps.Map{
		"protocol":  params.Protocol,
		"host":      "",
		"portRange": "",
	}

	// TODO 判断端口不能小于1
	// TODO 判断端口号不能大于65535

	digitRegexp := regexp.MustCompile(`^\d+$`)
	if digitRegexp.MatchString(params.Address) {
		addr["portRange"] = params.Address
	} else if strings.Contains(params.Address, ":") {
		index := strings.LastIndex(params.Address, ":")
		addr["host"] = strings.TrimSpace(params.Address[:index])
		port := strings.TrimSpace(params.Address[index+1:])
		if !digitRegexp.MatchString(port) {
			this.Fail("端口只能是一个数字")
		}
		addr["portRange"] = port
	} else {
		this.Fail("请输入正确的端口或者网络地址")
	}

	this.Data["address"] = addr
	this.Success()
}
