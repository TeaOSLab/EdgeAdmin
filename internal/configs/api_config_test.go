package configs

import (
	_ "github.com/iwind/TeaGo/bootstrap"
	"testing"
)

func TestLoadAPIConfig(t *testing.T) {
	config, err := LoadAPIConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config)
}

func TestAPIConfig_WriteFile(t *testing.T) {
	config := &APIConfig{
		RPC: struct {
			Endpoints []string `yaml:"endpoints"`
		}{},
		NodeId: "1",
		Secret: "2",
	}
	err := config.WriteFile("/tmp/api_config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
