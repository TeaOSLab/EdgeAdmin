package serverNames

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/serverutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

// 域名管理
type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.FirstMenu("index")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	server, _, isOk := serverutils.FindServer(this.Parent(), params.ServerId)
	if !isOk {
		return
	}

	serverNamesConfig := []*serverconfigs.ServerNameConfig{}
	if len(server.ServerNamesJSON) > 0 {
		err := json.Unmarshal(server.ServerNamesJSON, &serverNamesConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["serverNames"] = serverNamesConfig

	// DNS
	this.Data["dnsName"] = server.DnsName

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId    int64
	ServerNames string
	Must        *actions.Must
}) {
	// 记录日志
	defer this.CreateLog(oplogs.LevelInfo, "修改代理服务 %d 域名", params.ServerId)

	serverNames := []*serverconfigs.ServerNameConfig{}
	err := json.Unmarshal([]byte(params.ServerNames), &serverNames)
	if err != nil {
		this.Fail("域名解析失败：" + err.Error())
	}

	_, err = this.RPC().ServerRPC().UpdateServerNames(this.AdminContext(), &pb.UpdateServerNamesRequest{
		ServerId: params.ServerId,
		Config:   []byte(params.ServerNames),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
