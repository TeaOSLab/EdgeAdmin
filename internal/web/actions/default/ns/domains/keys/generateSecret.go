// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package keys

import (
	"encoding/base64"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/iwind/TeaGo/rands"
)

type GenerateSecretAction struct {
	actionutils.ParentAction
}

func (this *GenerateSecretAction) RunPost(params struct {
	SecretType string
}) {
	switch params.SecretType {
	case dnsconfigs.NSKeySecretTypeClear:
		this.Data["secret"] = rands.HexString(128)
	case dnsconfigs.NSKeySecretTypeBase64:
		this.Data["secret"] = base64.StdEncoding.EncodeToString([]byte(rands.HexString(128)))
	default:
		this.Data["secret"] = rands.HexString(128)
	}

	this.Success()
}
