package configloaders

import (
	"github.com/iwind/TeaGo/logs"
	"testing"
)

func TestLoadAdminModuleMapping(t *testing.T) {
	m, err := loadAdminModuleMapping()
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(m, t)
}
