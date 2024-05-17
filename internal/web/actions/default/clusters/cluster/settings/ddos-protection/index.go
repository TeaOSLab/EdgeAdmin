// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package ddosProtection

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/ddosconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	"net"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("ddosProtection")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	this.Data["clusterId"] = params.ClusterId

	protectionResp, err := this.RPC().NodeClusterRPC().FindNodeClusterDDoSProtection(this.AdminContext(), &pb.FindNodeClusterDDoSProtectionRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var ddosProtectionConfig = ddosconfigs.DefaultProtectionConfig()
	if len(protectionResp.DdosProtectionJSON) > 0 {
		err = json.Unmarshal(protectionResp.DdosProtectionJSON, ddosProtectionConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["config"] = ddosProtectionConfig
	this.Data["defaultConfigs"] = nodeconfigs.DefaultConfigs

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId          int64
	DdosProtectionJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.DDoSProtection_LogUpdateClusterDDoSProtection, params.ClusterId)

	var ddosProtectionConfig = &ddosconfigs.ProtectionConfig{}
	err := json.Unmarshal(params.DdosProtectionJSON, ddosProtectionConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	err = ddosProtectionConfig.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
	}

	// 校验参数
	if ddosProtectionConfig.TCP != nil {
		var tcpConfig = ddosProtectionConfig.TCP
		if tcpConfig.MaxConnectionsPerIP > 0 && tcpConfig.MaxConnectionsPerIP < nodeconfigs.DefaultTCPMinConnectionsPerIP {
			this.FailField("tcpMaxConnectionsPerIP", "TCP: 单IP TCP最大连接数不能小于"+types.String(nodeconfigs.DefaultTCPMinConnectionsPerIP))
		}

		if tcpConfig.NewConnectionsMinutelyRate > 0 && tcpConfig.NewConnectionsMinutelyRate < nodeconfigs.DefaultTCPNewConnectionsMinMinutelyRate {
			this.FailField("tcpNewConnectionsMinutelyRate", "TCP: 单IP连接速率不能小于"+types.String(nodeconfigs.DefaultTCPNewConnectionsMinMinutelyRate))
		}

		if tcpConfig.NewConnectionsSecondlyRate > 0 && tcpConfig.NewConnectionsSecondlyRate < nodeconfigs.DefaultTCPNewConnectionsMinSecondlyRate {
			this.FailField("tcpNewConnectionsSecondlyRate", "TCP: 单IP连接速率不能小于"+types.String(nodeconfigs.DefaultTCPNewConnectionsMinSecondlyRate))
		}

		// Port
		for _, portConfig := range tcpConfig.Ports {
			if portConfig.Port > 65535 {
				this.Fail("端口号" + types.String(portConfig.Port) + "不能大于65535")
			}
		}

		// IP
		for _, ipConfig := range tcpConfig.AllowIPList {
			if net.ParseIP(ipConfig.IP) == nil {
				this.Fail("白名单IP '" + ipConfig.IP + "' 格式错误")
			}
		}
	}

	_, err = this.RPC().NodeClusterRPC().UpdateNodeClusterDDoSProtection(this.AdminContext(), &pb.UpdateNodeClusterDDoSProtectionRequest{
		NodeClusterId:      params.ClusterId,
		DdosProtectionJSON: params.DdosProtectionJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Success()
}
