package grantutils

import "github.com/iwind/TeaGo/maps"

// 所有的认证类型
func AllGrantMethods() []maps.Map {
	return []maps.Map{
		{
			"name":  "用户名+密码",
			"value": "user",
		},
		{
			"name":  "私钥",
			"value": "privateKey",
		},
	}
}

// 获得对应的认证类型名称
func FindGrantMethodName(method string) string {
	for _, m := range AllGrantMethods() {
		if m.GetString("value") == method {
			return m.GetString("name")
		}
	}
	return ""
}
