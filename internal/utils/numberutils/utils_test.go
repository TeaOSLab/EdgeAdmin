// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package numberutils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"testing"
)

func TestFormatBytes(t *testing.T) {
	t.Log(numberutils.FormatBytes(1))
	t.Log(numberutils.FormatBytes(1000))
	t.Log(numberutils.FormatBytes(1_000_000))
	t.Log(numberutils.FormatBytes(1_000_000_000))
	t.Log(numberutils.FormatBytes(1_000_000_000_000))
	t.Log(numberutils.FormatBytes(1_000_000_000_000_000))
	t.Log(numberutils.FormatBytes(1_000_000_000_000_000_000))
	t.Log(numberutils.FormatBytes(9_000_000_000_000_000_000))
}

func TestFormatCount(t *testing.T) {
	t.Log(numberutils.FormatCount(1))
	t.Log(numberutils.FormatCount(1000))
	t.Log(numberutils.FormatCount(1500))
	t.Log(numberutils.FormatCount(1_000_000))
	t.Log(numberutils.FormatCount(1_500_000))
	t.Log(numberutils.FormatCount(1_000_000_000))
	t.Log(numberutils.FormatCount(1_500_000_000))
}

func TestFormatFloat(t *testing.T) {
	t.Log(numberutils.FormatFloat(1, 2))
	t.Log(numberutils.FormatFloat(100.23456, 2))
	t.Log(numberutils.FormatFloat(100.000023, 2))
	t.Log(numberutils.FormatFloat(100.012, 2))
}

func TestTrimZeroSuffix(t *testing.T) {
	for _, s := range []string{
		"1",
		"1.0000",
		"1.10",
		"100",
		"100.0000",
		"100.0",
		"100.0123",
		"100.0010",
		"100.000KB",
		"100.010MB",
	} {
		t.Log(s, "=>", numberutils.TrimZeroSuffix(s))
	}
}
