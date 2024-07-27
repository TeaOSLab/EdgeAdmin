// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cloud .

package utils_test

import (
	"testing"

	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/iwind/TeaGo/assert"
)

func TestValidateEmail(t *testing.T) {
	var a = assert.NewAssertion(t)
	a.IsTrue(utils.ValidateEmail("aaaa@gmail.com"))
	a.IsTrue(utils.ValidateEmail("a.b@gmail.com"))
	a.IsTrue(utils.ValidateEmail("a.b.c.d@gmail.com"))
	a.IsTrue(utils.ValidateEmail("aaaa@gmail.com.cn"))
	a.IsTrue(utils.ValidateEmail("hello.world.123@gmail.123.com"))
	a.IsTrue(utils.ValidateEmail("10000@qq.com"))
	a.IsFalse(utils.ValidateEmail("aaaa.@gmail.com"))
	a.IsFalse(utils.ValidateEmail("aaaa@gmail"))
	a.IsFalse(utils.ValidateEmail("aaaa@123"))
}
