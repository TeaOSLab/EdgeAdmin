// Copyright 2023 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package apinodeutils

var SharedManager = NewManager()

type Manager struct {
	upgraderMap map[int64]*Upgrader
}

func NewManager() *Manager {
	return &Manager{
		upgraderMap: map[int64]*Upgrader{},
	}
}

func (this *Manager) AddUpgrader(upgrader *Upgrader) {
	this.upgraderMap[upgrader.apiNodeId] = upgrader
}

func (this *Manager) FindUpgrader(apiNodeId int64) *Upgrader {
	return this.upgraderMap[apiNodeId]
}

func (this *Manager) RemoveUpgrader(upgrader *Upgrader) {
	if upgrader == nil {
		return
	}
	delete(this.upgraderMap, upgrader.apiNodeId)
}
