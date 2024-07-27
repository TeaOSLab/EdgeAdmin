package configloaders

import (
	"testing"

	"github.com/iwind/TeaGo/logs"
)

func TestLoadAdminModuleMapping(t *testing.T) {
	m, err := loadAdminModuleMapping()
	if err != nil {
		t.Fatal(err)
	}
	logs.PrintAsJSON(m, t)
}
