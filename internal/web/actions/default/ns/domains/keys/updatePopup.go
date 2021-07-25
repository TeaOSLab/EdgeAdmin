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

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	KeyId int64
}) {
	keyResp, err := this.RPC().NSKeyRPC().FindEnabledNSKey(this.AdminContext(), &pb.FindEnabledNSKeyRequest{NsKeyId: params.KeyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var key = keyResp.NsKey
	if key == nil {
		return
	}

	this.Data["key"] = maps.Map{
		"id":         key.Id,
		"name":       key.Name,
		"algo":       key.Algo,
		"secret":     key.Secret,
		"secretType": key.SecretType,
		"isOn":       key.IsOn,
	}

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

func (this *UpdatePopupAction) RunPost(params struct {
	KeyId      int64
	Name       string
	Algo       string
	Secret     string
	SecretType string
	IsOn       bool

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	this.CreateLogInfo("修改DNS密钥 %d", params.KeyId)

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

	_, err := this.RPC().NSKeyRPC().UpdateNSKey(this.AdminContext(), &pb.UpdateNSKeyRequest{
		NsKeyId:    params.KeyId,
		Name:       params.Name,
		Algo:       params.Algo,
		Secret:     params.Secret,
		SecretType: params.SecretType,
		IsOn:       params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
