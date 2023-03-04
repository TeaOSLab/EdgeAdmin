// Copyright 2023 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package apinodeutils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/apinodeutils"
	_ "github.com/iwind/TeaGo/bootstrap"
	"testing"
)

func TestUpgrader_CanUpgrade(t *testing.T) {
	var upgrader = apinodeutils.NewUpgrader()
	t.Log(upgrader.CanUpgrade("0.6.3"))
}

func TestUpgrader_Upgrade(t *testing.T) {
	var upgrader = apinodeutils.NewUpgrader()
	err := upgrader.Upgrade(1)
	if err != nil {
		t.Fatal(err)
	}
}
