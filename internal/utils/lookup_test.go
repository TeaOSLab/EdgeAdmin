// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package utils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"testing"
)

func TestLookupCNAME(t *testing.T) {
	for _, domain := range []string{"www.yun4s.cn", "example.com"} {
		result, err := utils.LookupCNAME(domain)
		t.Log(domain, "=>", result, err)
	}
}
