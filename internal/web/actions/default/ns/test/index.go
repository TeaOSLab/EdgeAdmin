// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"github.com/miekg/dns"
	"net"
	"regexp"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	// 集群列表
	clustersResp, err := this.RPC().NSClusterRPC().FindAllEnabledNSClusters(this.AdminContext(), &pb.FindAllEnabledNSClustersRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var clusterMaps = []maps.Map{}
	for _, cluster := range clustersResp.NsClusters {
		if !cluster.IsOn {
			continue
		}

		countNodesResp, err := this.RPC().NSNodeRPC().CountAllEnabledNSNodesMatch(this.AdminContext(), &pb.CountAllEnabledNSNodesMatchRequest{
			NsClusterId:  cluster.Id,
			InstallState: 0,
			ActiveState:  0,
			Keyword:      "",
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var countNodes = countNodesResp.Count
		if countNodes <= 0 {
			continue
		}

		clusterMaps = append(clusterMaps, maps.Map{
			"id":         cluster.Id,
			"name":       cluster.Name,
			"countNodes": countNodes,
		})
	}
	this.Data["clusters"] = clusterMaps

	// 记录类型
	this.Data["recordTypes"] = dnsconfigs.FindAllRecordTypeDefinitions()

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	NodeId   int64
	Domain   string
	Type     string
	Ip       string
	ClientIP string

	Must *actions.Must
}) {
	nodeResp, err := this.RPC().NSNodeRPC().FindEnabledNSNode(this.AdminContext(), &pb.FindEnabledNSNodeRequest{NsNodeId: params.NodeId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var node = nodeResp.NsNode
	if node == nil {
		this.Fail("找不到要测试的节点")
	}

	var isOk = false
	var errMsg string
	var isNetError = false
	var result string

	defer func() {
		this.Data["isOk"] = isOk
		this.Data["err"] = errMsg
		this.Data["isNetErr"] = isNetError
		this.Data["result"] = result
		this.Success()
	}()

	if !domainutils.ValidateDomainFormat(params.Domain) {
		errMsg = "域名格式错误"
		return
	}

	recordType, ok := dns.StringToType[params.Type]
	if !ok {
		errMsg = "不支持此记录类型"
		return
	}

	if len(params.ClientIP) > 0 && net.ParseIP(params.ClientIP) == nil {
		errMsg = "客户端IP格式不正确"
		return
	}

	var optionId int64
	if len(params.ClientIP) > 0 {
		optionResp, err := this.RPC().NSQuestionOptionRPC().CreateNSQuestionOption(this.AdminContext(), &pb.CreateNSQuestionOptionRequest{
			Name:       "setRemoteAddr",
			ValuesJSON: maps.Map{"ip": params.ClientIP}.AsJSON(),
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		optionId = optionResp.NsQuestionOptionId
		defer func() {
			_, err = this.RPC().NSQuestionOptionRPC().DeleteNSQuestionOption(this.AdminContext(), &pb.DeleteNSQuestionOptionRequest{NsQuestionOptionId: optionId})
			if err != nil {
				this.ErrorPage(err)
			}
		}()
	}

	c := new(dns.Client)
	m := new(dns.Msg)
	var domain = params.Domain + "."
	if optionId > 0 {
		domain = "$" + types.String(optionId) + "-" + domain
	}
	m.SetQuestion(domain, recordType)
	r, _, err := c.Exchange(m, params.Ip+":53")
	if err != nil {
		errMsg = "解析过程中出错：" + err.Error()

		// 是否为网络错误
		if regexp.MustCompile(`timeout|connect`).MatchString(err.Error()) {
			isNetError = true
		}

		return
	}
	result = r.String()
	result = regexp.MustCompile(`\$\d+-`).ReplaceAllString(result, "")
	isOk = true
}
