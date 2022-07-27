package servers

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
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
	ServerType   string
	Protocol     string
	From         string
	SupportRange bool
}) {
	this.Data["from"] = params.From

	protocols := serverconfigs.FindAllServerProtocolsForType(params.ServerType)
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

	this.Data["supportRange"] = params.SupportRange

	this.Show()
}

func (this *AddPortPopupAction) RunPost(params struct {
	SupportRange bool

	Protocol string
	Address  string

	Must *actions.Must
}) {
	// 校验地址
	addr := maps.Map{
		"protocol":  params.Protocol,
		"host":      "",
		"portRange": "",
		"minPort":   0,
		"maxPort":   0,
	}

	var portRegexp = regexp.MustCompile(`^\d+$`)
	if portRegexp.MatchString(params.Address) { // 单个端口
		addr["portRange"] = this.checkPort(params.Address)
	} else if params.SupportRange && regexp.MustCompile(`^\d+\s*-\s*\d+$`).MatchString(params.Address) { // Port1-Port2
		addr["portRange"], addr["minPort"], addr["maxPort"] = this.checkPortRange(params.Address)
	} else if strings.Contains(params.Address, ":") { // IP:Port
		index := strings.LastIndex(params.Address, ":")
		addr["host"] = strings.TrimSpace(params.Address[:index])
		port := strings.TrimSpace(params.Address[index+1:])
		if portRegexp.MatchString(port) {
			addr["portRange"] = this.checkPort(port)
		} else if params.SupportRange && regexp.MustCompile(`^\d+\s*-\s*\d+$`).MatchString(port) { // Port1-Port2
			addr["portRange"], addr["minPort"], addr["maxPort"] = this.checkPortRange(port)
		} else {
			this.FailField("address", "请输入正确的端口或者网络地址")
		}
	} else {
		this.FailField("address", "请输入正确的端口或者网络地址")
	}

	this.Data["address"] = addr
	this.Success()
}

func (this *AddPortPopupAction) checkPort(port string) (portRange string) {
	var intPort = types.Int(port)
	if intPort < 1 {
		this.FailField("address", "端口号不能小于1")
	}
	if intPort > 65535 {
		this.FailField("address", "端口号不能大于65535")
	}
	return port
}

func (this *AddPortPopupAction) checkPortRange(port string) (portRange string, minPort int, maxPort int) {
	var pieces = strings.Split(port, "-")
	var piece1 = strings.TrimSpace(pieces[0])
	var piece2 = strings.TrimSpace(pieces[1])
	var port1 = types.Int(piece1)
	var port2 = types.Int(piece2)

	if port1 < 1 {
		this.FailField("address", "端口号不能小于1")
	}
	if port1 > 65535 {
		this.FailField("address", "端口号不能大于65535")
	}

	if port2 < 1 {
		this.FailField("address", "端口号不能小于1")
	}
	if port2 > 65535 {
		this.FailField("address", "端口号不能大于65535")
	}

	if port1 > port2 {
		port1, port2 = port2, port1
	}

	return types.String(port1) + "-" + types.String(port2), port1, port2
}
