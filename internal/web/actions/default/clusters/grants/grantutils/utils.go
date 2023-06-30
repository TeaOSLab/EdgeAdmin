package grantutils

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/langs"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/iwind/TeaGo/maps"
)

// AllGrantMethods 所有的认证类型
func AllGrantMethods(langCode langs.LangCode) []maps.Map {
	return []maps.Map{
		{
			"name":  langs.Message(langCode, codes.NodeGrant_MethodUserPassword),
			"value": "user",
		},
		{
			"name":  langs.Message(langCode, codes.NodeGrant_MethodPrivateKey),
			"value": "privateKey",
		},
	}
}

// FindGrantMethodName 获得对应的认证类型名称
func FindGrantMethodName(method string, langCode langs.LangCode) string {
	for _, m := range AllGrantMethods(langCode) {
		if m.GetString("value") == method {
			return m.GetString("name")
		}
	}
	return ""
}
