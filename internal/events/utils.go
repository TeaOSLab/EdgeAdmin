package events

import "sync"

var eventsMap = map[string][]func(){} // event => []callbacks
var locker = sync.Mutex{}

// On 增加事件回调
func On(event string, callback func()) {
	locker.Lock()
	defer locker.Unlock()

	var callbacks = eventsMap[event]
	callbacks = append(callbacks, callback)
	eventsMap[event] = callbacks
}

// Notify 通知事件
func Notify(event string) {
	locker.Lock()
	var callbacks = eventsMap[event]
	locker.Unlock()

	for _, callback := range callbacks {
		callback()
	}
}
