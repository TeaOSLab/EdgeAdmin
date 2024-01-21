// Copyright 2024 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package nodeutils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/clusters/cluster/node/nodeutils"
	_ "github.com/iwind/TeaGo/bootstrap"
	"testing"
)

func TestInstallLocalNode(t *testing.T) {
	err := nodeutils.InstallLocalNode()
	if err != nil {
		t.Fatal(err)
	}
}
