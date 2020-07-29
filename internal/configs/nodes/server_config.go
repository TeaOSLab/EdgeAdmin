package nodes

import "encoding/json"

type ServerConfig struct {
	Id   string `json:"id" yaml:"id"`
	IsOn bool   `json:"isOn" yaml:"isOn"`
	Name string `json:"name" yaml:"name"`
}

func (this *ServerConfig) AsJSON() ([]byte, error) {
	return json.Marshal(this)
}
