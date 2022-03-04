// Copyright 2022 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package utils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
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
