package configloaders

import (
	"testing"
	"time"

	_ "github.com/iwind/TeaGo/bootstrap"
)

func TestLoadUIConfig(t *testing.T) {
	for i := 0; i < 10; i++ {
		before := time.Now()
		config, err := LoadAdminUIConfig()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(time.Since(before).Seconds()*1000, "ms")
		t.Logf("%p", config)
	}
}

func TestLoadUIConfig2(t *testing.T) {
	for i := 0; i < 10; i++ {
		config, err := LoadAdminUIConfig()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(config)
	}
}
