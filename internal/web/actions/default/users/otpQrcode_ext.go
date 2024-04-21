// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .
//go:build !plus

package users

import "github.com/TeaOSLab/EdgeAdmin/internal/configloaders"

func (this *OtpQrcodeAction) findProductName() (string, error) {
	uiConfig, err := configloaders.LoadAdminUIConfig()
	if err != nil {
		return "", err
	}
	var productName = uiConfig.ProductName
	if len(productName) == 0 {
		productName = "GoEdge用户"
	}
	return productName, nil
}
