// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package utils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestJSONClone(t *testing.T) {
	type A struct {
		B int    `json:"b"`
		C string `json:"c"`
	}

	var a = &A{B: 123, C: "456"}

	for i := 0; i < 5; i++ {
		c, err := utils.JSONClone(a)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%p, %#v", c, c)
	}
}


func TestJSONIsNull(t *testing.T) {
	var a = assert.NewAssertion(t)
	a.IsTrue(utils.JSONIsNull(nil))
	a.IsTrue(utils.JSONIsNull([]byte{}))
	a.IsTrue(utils.JSONIsNull([]byte("null")))
	a.IsFalse(utils.JSONIsNull([]byte{1, 2, 3}))
}