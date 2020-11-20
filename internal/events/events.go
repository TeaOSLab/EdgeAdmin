package events

type Event = string

const (
	EventStart Event = "start" // start loading
	EventQuit  Event = "quit"  // quit node gracefully

	EventSecurityConfigChanged Event = "securityConfigChanged" // 安全设置变更
)
