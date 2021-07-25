// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package keys

import (
	"encoding/base64"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct {
	DomainId int64
}) {
	this.Data["domainId"] = params.DomainId

	// 所有算法
	var algorithmMaps = []maps.Map{}
	for _, algo := range dnsconfigs.FindAllKeyAlgorithmTypes() {
		algorithmMaps = append(algorithmMaps, maps.Map{
			"name": algo.Name,
			"code": algo.Code,
		})
	}
	this.Data["algorithms"] = algorithmMaps

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	DomainId   int64
	Name       string
	Algo       string
	Secret     string
	SecretType string

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	var keyId int64 = 0
	defer func() {
		this.CreateLogInfo("创建DNS密钥 %d", keyId)
	}()

	params.Must.
		Field("name", params.Name).
		Require("请输入密钥名称").
		Field("algo", params.Algo).
		Require("请选择算法").
		Field("secret", params.Secret).
		Require("请输入密码")

	// 校验密码
	if params.SecretType == dnsconfigs.NSKeySecretTypeBase64 {
		_, err := base64.StdEncoding.DecodeString(params.Secret)
		if err != nil {
			this.FailField("secret", "请输入BASE64格式的密码或者选择明文")
		}
	}

	createResp, err := this.RPC().NSKeyRPC().CreateNSKey(this.AdminContext(), &pb.CreateNSKeyRequest{
		NsDomainId: params.DomainId,
		NsZoneId:   0,
		Name:       params.Name,
		Algo:       params.Algo,
		Secret:     params.Secret,
		SecretType: params.SecretType,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	keyId = createResp.NsKeyId

	this.Success()
}
