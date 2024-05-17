// Copyright 2021 GoEdge CDN goedge.cdn@gmail.com. All rights reserved.

package numberutils_test

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils/numberutils"
	"github.com/iwind/TeaGo/assert"
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
	t.Log(numberutils.FormatFloat(123.012, 2))
	t.Log(numberutils.FormatFloat(1234.012, 2))
	t.Log(numberutils.FormatFloat(12345.012, 2))
	t.Log(numberutils.FormatFloat(123456.012, 2))
	t.Log(numberutils.FormatFloat(1234567.012, 2))
	t.Log(numberutils.FormatFloat(12345678.012, 2))
	t.Log(numberutils.FormatFloat(123456789.012, 2))
	t.Log(numberutils.FormatFloat(1234567890.012, 2))
	t.Log(numberutils.FormatFloat(123, 2))
	t.Log(numberutils.FormatFloat(1234, 2))
	t.Log(numberutils.FormatFloat(1234.00001, 4))
	t.Log(numberutils.FormatFloat(1234.56700, 4))
	t.Log(numberutils.FormatFloat(-1234.56700, 2))
	t.Log(numberutils.FormatFloat(-221745.12, 2))
}

func TestFormatFloat2(t *testing.T) {
	t.Log(numberutils.FormatFloat2(0))
	t.Log(numberutils.FormatFloat2(0.0))
	t.Log(numberutils.FormatFloat2(1.23456))
	t.Log(numberutils.FormatFloat2(1.0))
}

func TestPadFloatZero(t *testing.T) {
	var a = assert.NewAssertion(t)
	a.IsTrue(numberutils.PadFloatZero("1", 0) == "1")
	a.IsTrue(numberutils.PadFloatZero("1", 2) == "1.00")
	a.IsTrue(numberutils.PadFloatZero("1.1", 2) == "1.10")
	a.IsTrue(numberutils.PadFloatZero("1.12", 2) == "1.12")
	a.IsTrue(numberutils.PadFloatZero("1.123", 2) == "1.123")
	a.IsTrue(numberutils.PadFloatZero("10000.123", 2) == "10000.123")
	a.IsTrue(numberutils.PadFloatZero("", 2) == "0.00")
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
