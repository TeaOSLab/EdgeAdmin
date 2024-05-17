// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package fastcgi

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"net"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Address         string
	ParamsJSON      []byte
	ReadTimeout     int64
	PoolSize        int32
	PathInfoPattern string
	IsOn            bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var fastcgiId = int64(0)
	defer func() {
		this.CreateLogInfo(codes.HTTPFastcgi_LogCreateHTTPFastcgi, fastcgiId)
	}()

	params.Must.
		Field("address", params.Address).
		Require("请输入Fastcgi地址")

	_, _, err := net.SplitHostPort(params.Address)
	if err != nil {
		this.FailField("address", "请输入正确的Fastcgi地址")
	}

	readTimeoutJSON, err := json.Marshal(&shared.TimeDuration{
		Count: params.ReadTimeout,
		Unit:  "second",
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	createResp, err := this.RPC().HTTPFastcgiRPC().CreateHTTPFastcgi(this.AdminContext(), &pb.CreateHTTPFastcgiRequest{
		IsOn:            params.IsOn,
		Address:         params.Address,
		ParamsJSON:      params.ParamsJSON,
		ReadTimeoutJSON: readTimeoutJSON,
		ConnTimeoutJSON: nil, // TODO 将来支持
		PoolSize:        params.PoolSize,
		PathInfoPattern: params.PathInfoPattern,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	fastcgiId = createResp.HttpFastcgiId

	configResp, err := this.RPC().HTTPFastcgiRPC().FindEnabledHTTPFastcgiConfig(this.AdminContext(), &pb.FindEnabledHTTPFastcgiConfigRequest{HttpFastcgiId: fastcgiId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	configJSON := configResp.HttpFastcgiJSON
	config := &serverconfigs.HTTPFastcgiConfig{}
	err = json.Unmarshal(configJSON, config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["fastcgi"] = config

	this.Success()
}
