// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package sizes_test

import (
	"testing"

	"github.com/TeaOSLab/EdgeAdmin/internal/utils/sizes"
	"github.com/iwind/TeaGo/assert"
)

func TestSizes(t *testing.T) {
	var a = assert.NewAssertion(t)
	a.IsTrue(sizes.K == 1024)
	a.IsTrue(sizes.M == 1024*1024)
	a.IsTrue(sizes.G == 1024*1024*1024)
	a.IsTrue(sizes.T == 1024*1024*1024*1024)
}
