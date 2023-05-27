package actionutils

import (
	"github.com/iwind/TeaGo/actions"
)

type TabItem struct {
	Name       string `json:"name"`
	SubName    string `json:"subName"`
	URL        string `json:"url"`
	Icon       string `json:"icon"`
	IsActive   bool   `json:"isActive"`
	IsRight    bool   `json:"isRight"`
	IsTitle    bool   `json:"isTitle"`
	IsDisabled bool   `json:"isDisabled"`
}

// Tabbar Tabbar定义
type Tabbar struct {
	items []*TabItem
}

// NewTabbar 获取新对象
func NewTabbar() *Tabbar {
	return &Tabbar{
		items: []*TabItem{},
	}
}

// Add 添加菜单项
func (this *Tabbar) Add(name string, subName string, url string, icon string, active bool) *TabItem {
	var m = &TabItem{
		Name:       name,
		SubName:    subName,
		URL:        url,
		Icon:       icon,
		IsActive:   active,
		IsRight:    false,
		IsTitle:    false,
		IsDisabled: false,
	}
	this.items = append(this.items, m)
	return m
}

// Items 取得所有的Items
func (this *Tabbar) Items() []*TabItem {
	return this.items
}

// SetTabbar 设置子菜单
func SetTabbar(action actions.ActionWrapper, tabbar *Tabbar) {
	action.Object().Data["teaTabbar"] = tabbar.Items()
}
