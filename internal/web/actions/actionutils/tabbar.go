package actionutils

import (
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

// Tabbar Tabbar定义
type Tabbar struct {
	items []maps.Map
}

// NewTabbar 获取新对象
func NewTabbar() *Tabbar {
	return &Tabbar{
		items: []maps.Map{},
	}
}

// Add 添加菜单项
func (this *Tabbar) Add(name string, subName string, url string, icon string, active bool) maps.Map {
	m := maps.Map{
		"name":    name,
		"subName": subName,
		"url":     url,
		"icon":    icon,
		"active":  active,
		"right":   false,
		"isTitle": false,
	}
	this.items = append(this.items, m)
	return m
}

// Items 取得所有的Items
func (this *Tabbar) Items() []maps.Map {
	return this.items
}

// SetTabbar 设置子菜单
func SetTabbar(action actions.ActionWrapper, tabbar *Tabbar) {
	action.Object().Data["teaTabbar"] = tabbar.Items()
}
