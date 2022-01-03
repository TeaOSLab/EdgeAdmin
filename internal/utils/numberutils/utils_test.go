// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package numberutils

import "testing"

func TestFormatBytes(t *testing.T) {
	t.Log(FormatBytes(1))
	t.Log(FormatBytes(1000))
	t.Log(FormatBytes(1_000_000))
	t.Log(FormatBytes(1_000_000_000))
	t.Log(FormatBytes(1_000_000_000_000))
	t.Log(FormatBytes(1_000_000_000_000_000))
	t.Log(FormatBytes(1_000_000_000_000_000_000))
	t.Log(FormatBytes(9_000_000_000_000_000_000))
}

func TestFormatCount(t *testing.T) {
	t.Log(FormatCount(1))
	t.Log(FormatCount(1000))
	t.Log(FormatCount(1500))
	t.Log(FormatCount(1_000_000))
	t.Log(FormatCount(1_500_000))
	t.Log(FormatCount(1_000_000_000))
	t.Log(FormatCount(1_500_000_000))
}