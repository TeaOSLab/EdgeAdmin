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
