// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package gen

import "testing"

func TestGenerate(t *testing.T) {
	err := Generate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("ok")
}
