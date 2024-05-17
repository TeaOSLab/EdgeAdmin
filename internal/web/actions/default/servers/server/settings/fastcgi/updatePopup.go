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

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	FastcgiId int64
}) {
	configResp, err := this.RPC().HTTPFastcgiRPC().FindEnabledHTTPFastcgiConfig(this.AdminContext(), &pb.FindEnabledHTTPFastcgiConfigRequest{HttpFastcgiId: params.FastcgiId})
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

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	FastcgiId       int64
	Address         string
	ParamsJSON      []byte
	ReadTimeout     int64
	PoolSize        int32
	PathInfoPattern string
	IsOn            bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.HTTPFastcgi_LogUpdateHTTPFastcgi, params.FastcgiId)

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

	_, err = this.RPC().HTTPFastcgiRPC().UpdateHTTPFastcgi(this.AdminContext(), &pb.UpdateHTTPFastcgiRequest{
		HttpFastcgiId:   params.FastcgiId,
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

	configResp, err := this.RPC().HTTPFastcgiRPC().FindEnabledHTTPFastcgiConfig(this.AdminContext(), &pb.FindEnabledHTTPFastcgiConfigRequest{HttpFastcgiId: params.FastcgiId})
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
