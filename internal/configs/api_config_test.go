package configs

import (
	_ "github.com/iwind/TeaGo/bootstrap"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestLoadAPIConfig(t *testing.T) {
	config, err := LoadAPIConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(config)

	configData, err := yaml.Marshal(config)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(configData))
}

func TestAPIConfig_WriteFile(t *testing.T) {
	var config = &APIConfig{}
	err := config.WriteFile("/tmp/api_config.yaml")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
