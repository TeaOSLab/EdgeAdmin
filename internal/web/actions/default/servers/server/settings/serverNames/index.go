package serverNames

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"strings"
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
	serverNamesResp, err := this.RPC().ServerRPC().FindServerNames(this.AdminContext(), &pb.FindServerNamesRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	serverNamesConfig := []*serverconfigs.ServerNameConfig{}
	this.Data["isAuditing"] = serverNamesResp.IsAuditing
	this.Data["auditingResult"] = maps.Map{
		"isOk": true,
	}
	if serverNamesResp.IsAuditing {
		serverNamesResp.ServerNamesJSON = serverNamesResp.AuditingServerNamesJSON
	} else if serverNamesResp.AuditingResult != nil {
		if !serverNamesResp.AuditingResult.IsOk {
			serverNamesResp.ServerNamesJSON = serverNamesResp.AuditingServerNamesJSON
		}

		this.Data["auditingResult"] = maps.Map{
			"isOk":        serverNamesResp.AuditingResult.IsOk,
			"reason":      serverNamesResp.AuditingResult.Reason,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", serverNamesResp.AuditingResult.CreatedAt),
		}
	}
	if len(serverNamesResp.ServerNamesJSON) > 0 {
		err := json.Unmarshal(serverNamesResp.ServerNamesJSON, &serverNamesConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["serverNames"] = serverNamesConfig

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ServerId    int64
	ServerNames string
	Must        *actions.Must
	CSRF        *actionutils.CSRF
}) {
	// 记录日志
	defer this.CreateLog(oplogs.LevelInfo, "修改代理服务 %d 域名", params.ServerId)

	serverNames := []*serverconfigs.ServerNameConfig{}
	err := json.Unmarshal([]byte(params.ServerNames), &serverNames)
	if err != nil {
		this.Fail("域名解析失败：" + err.Error())
	}

	serverResp, err := this.RPC().ServerRPC().FindEnabledUserServerBasic(this.AdminContext(), &pb.FindEnabledUserServerBasicRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if serverResp.Server == nil || serverResp.Server.NodeCluster == nil {
		this.NotFound("server", params.ServerId)
		return
	}
	clusterId := serverResp.Server.NodeCluster.Id

	// 检查域名是否已经存在
	allServerNames := serverconfigs.PlainServerNames(serverNames)
	if len(allServerNames) > 0 {
		dupResp, err := this.RPC().ServerRPC().CheckServerNameDuplicationInNodeCluster(this.AdminContext(), &pb.CheckServerNameDuplicationInNodeClusterRequest{
			ServerNames:     allServerNames,
			NodeClusterId:   clusterId,
			ExcludeServerId: params.ServerId,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if len(dupResp.DuplicatedServerNames) > 0 {
			this.Fail("域名 " + strings.Join(dupResp.DuplicatedServerNames, ", ") + " 已经被其他服务所占用，不能重复使用")
		}
	}

	_, err = this.RPC().ServerRPC().UpdateServerNames(this.AdminContext(), &pb.UpdateServerNamesRequest{
		ServerId:        params.ServerId,
		ServerNamesJSON: []byte(params.ServerNames),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
