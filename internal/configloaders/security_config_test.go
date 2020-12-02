package configloaders

import (
	_ "github.com/iwind/TeaGo/bootstrap"
	"testing"
	"time"
)

func TestLoadSecurityConfig(t *testing.T) {
	for i := 0; i < 10; i++ {
		before := time.Now()
		config, err := LoadSecurityConfig()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(time.Since(before).Seconds()*1000, "ms")
		t.Logf("%p", config)
	}
}

func TestLoadSecurityConfig2(t *testing.T) {
	for i := 0; i < 10; i++ {
		config, err := LoadSecurityConfig()
		if err != nil {
			t.Fatal(err)
		}
		if i == 0 {
			config.Frame = "DENY"
		}
		t.Log(config.Frame)
	}
}
