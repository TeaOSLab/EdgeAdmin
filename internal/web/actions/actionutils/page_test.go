package actionutils

import (
	"net/url"
	"testing"
)

func TestNewActionPage(t *testing.T) {
	page := &Page{
		Current: 3,
		Total:   105,
		Size:    20,
		Path:    "/hello",
		Query: url.Values{
			"a":    []string{"b"},
			"c":    []string{"d"},
			"page": []string{"3"},
		},
	}
	page.calculate()
	t.Log(page.AsHTML())
	//logs.PrintAsJSON(page, t)
}

func TestNewActionPage2(t *testing.T) {
	page := &Page{
		Current: 3,
		Total:   105,
		Size:    10,
		Path:    "/hello",
		Query: url.Values{
			"a":    []string{"b"},
			"c":    []string{"d"},
			"page": []string{"3"},
		},
	}
	page.calculate()
	t.Log(page.AsHTML())
	//logs.PrintAsJSON(page, t)
}
