// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package ipbox

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/cmd"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Data["types"] = serverconfigs.FindAllAccessLogStorageTypes()
	this.Data["syslogPriorities"] = serverconfigs.AccessLogSyslogStoragePriorities

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name string
	Type string

	// file
	FilePath       string
	FileAutoCreate bool

	// es
	EsEndpoint    string
	EsIndex       string
	EsMappingType string
	EsUsername    string
	EsPassword    string

	// mysql
	MysqlHost     string
	MysqlPort     int
	MysqlUsername string
	MysqlPassword string
	MysqlDatabase string
	MysqlTable    string
	MysqlLogField string

	// tcp
	TcpNetwork string
	TcpAddr    string

	// syslog
	SyslogProtocol   string
	SyslogServerAddr string
	SyslogServerPort int
	SyslogSocket     string
	SyslogTag        string
	SyslogPriority   int

	// command
	CommandCommand string
	CommandArgs    string
	CommandDir     string

	IsPublic bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var policyId int64 = 0
	defer func() {
		this.CreateLogInfo("创建访问日志策略 %d", policyId)
	}()

	params.Must.
		Field("name", params.Name).
		Require("请输入日志策略的名称").
		Field("type", params.Type).
		Require("请选择存储类型")

	var options interface{} = nil
	switch params.Type {
	case serverconfigs.AccessLogStorageTypeFile:
		params.Must.
			Field("filePath", params.FilePath).
			Require("请输入日志文件路径")

		storage := new(serverconfigs.AccessLogFileStorageConfig)
		storage.Path = params.FilePath
		storage.AutoCreate = params.FileAutoCreate
		options = storage
	case serverconfigs.AccessLogStorageTypeES:
		params.Must.
			Field("esEndpoint", params.EsEndpoint).
			Require("请输入Endpoint").
			Field("esIndex", params.EsIndex).
			Require("请输入Index名称").
			Field("esMappingType", params.EsMappingType).
			Require("请输入Mapping名称")

		storage := new(serverconfigs.AccessLogESStorageConfig)
		storage.Endpoint = params.EsEndpoint
		storage.Index = params.EsIndex
		storage.MappingType = params.EsMappingType
		storage.Username = params.EsUsername
		storage.Password = params.EsPassword
		options = storage
	case serverconfigs.AccessLogStorageTypeTCP:
		params.Must.
			Field("tcpNetwork", params.TcpNetwork).
			Require("请选择网络协议").
			Field("tcpAddr", params.TcpAddr).
			Require("请输入网络地址")

		storage := new(serverconfigs.AccessLogTCPStorageConfig)
		storage.Network = params.TcpNetwork
		storage.Addr = params.TcpAddr
		options = storage
	case serverconfigs.AccessLogStorageTypeSyslog:
		switch params.SyslogProtocol {
		case serverconfigs.AccessLogSyslogStorageProtocolTCP, serverconfigs.AccessLogSyslogStorageProtocolUDP:
			params.Must.
				Field("syslogServerAddr", params.SyslogServerAddr).
				Require("请输入网络地址")
		case serverconfigs.AccessLogSyslogStorageProtocolSocket:
			params.Must.
				Field("syslogSocket", params.SyslogSocket).
				Require("请输入Socket路径")
		}

		storage := new(serverconfigs.AccessLogSyslogStorageConfig)
		storage.Protocol = params.SyslogProtocol
		storage.ServerAddr = params.SyslogServerAddr
		storage.ServerPort = params.SyslogServerPort
		storage.Socket = params.SyslogSocket
		storage.Tag = params.SyslogTag
		storage.Priority = params.SyslogPriority
		options = storage
	case serverconfigs.AccessLogStorageTypeCommand:
		params.Must.
			Field("commandCommand", params.CommandCommand).
			Require("请输入可执行命令")

		storage := new(serverconfigs.AccessLogCommandStorageConfig)
		storage.Command = params.CommandCommand
		storage.Args = cmd.ParseArgs(params.CommandArgs)
		storage.Dir = params.CommandDir
		options = storage
	}

	if options == nil {
		this.Fail("找不到选择的存储类型")
	}

	optionsJSON, err := json.Marshal(options)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	createResp, err := this.RPC().HTTPAccessLogPolicyRPC().CreateHTTPAccessLogPolicy(this.AdminContext(), &pb.CreateHTTPAccessLogPolicyRequest{
		Name:        params.Name,
		Type:        params.Type,
		OptionsJSON: optionsJSON,
		CondsJSON:   nil, // TODO
		IsPublic:    params.IsPublic,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyId = createResp.HttpAccessLogPolicyId

	this.Success()
}
