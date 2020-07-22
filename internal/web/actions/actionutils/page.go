package actionutils

import (
	"github.com/iwind/TeaGo/actions"
	"math"
	"net/url"
	"strconv"
	"strings"
)

type Page struct {
	Offset  int // 开始位置
	Size    int // 每页显示数量
	Current int // 当前页码
	Max     int // 最大页码
	Total   int // 总数量

	Path  string
	Query url.Values
}

func NewActionPage(actionPtr actions.ActionWrapper, total int, size int) *Page {
	action := actionPtr.Object()
	currentPage := action.ParamInt("page")

	paramSize := action.ParamInt("pageSize")
	if paramSize > 0 {
		size = paramSize
	}

	if size <= 0 {
		size = 10
	}

	page := &Page{
		Current: currentPage,
		Total:   total,
		Size:    size,
		Path:    action.Request.URL.Path,
		Query:   action.Request.URL.Query(),
	}
	page.calculate()
	return page
}

func (this *Page) calculate() {
	if this.Current < 1 {
		this.Current = 1
	}
	if this.Size <= 0 {
		this.Size = 10
	}

	this.Offset = this.Size * (this.Current - 1)
	this.Max = int(math.Ceil(float64(this.Total) / float64(this.Size)))
}

func (this *Page) AsHTML() string {
	if this.Total <= this.Size {
		return ""
	}

	result := []string{}

	// 首页
	if this.Max > 0 {
		result = append(result, `<a href="`+this.composeURL(1)+`">首页</a>`)
	} else {
		result = append(result, `<a>首页</a>`)
	}

	// 上一页
	if this.Current <= 1 {
		result = append(result, `<a>上一页</a>`)
	} else {
		result = append(result, `<a href="`+this.composeURL(this.Current-1)+`">上一页</a>`)
	}

	// 中间页数
	before5 := this.max(this.Current-5, 1)
	after5 := this.min(before5+9, this.Max)

	if before5 > 1 {
		result = append(result, `<a>...</a>`)
	}

	for i := before5; i <= after5; i++ {
		if i == this.Current {
			result = append(result, `<a href="`+this.composeURL(i)+`" class="active">`+strconv.Itoa(i)+`</a>`)
		} else {
			result = append(result, `<a href="`+this.composeURL(i)+`">`+strconv.Itoa(i)+`</a>`)
		}
	}

	if after5 < this.Max {
		result = append(result, `<a>...</a>`)
	}

	// 下一页
	if this.Current >= this.Max {
		result = append(result, "<a>下一页</a>")
	} else {
		result = append(result, `<a href="`+this.composeURL(this.Current+1)+`">下一页</a>`)
	}

	// 尾页
	if this.Max > 0 {
		result = append(result, `<a href="`+this.composeURL(this.Max)+`">尾页</a>`)
	} else {
		result = append(result, `<a>尾页</a>`)
	}

	// 每页数
	result = append(result, `<select class="ui dropdown" style="height:34px;padding-top:0;padding-bottom:0;margin-left:1em;color:#666" onchange="ChangePageSize(this.value)">
	<option value="10">[每页]</option>`+this.renderSizeOption(10)+
		this.renderSizeOption(20)+
		this.renderSizeOption(30)+
		this.renderSizeOption(40)+
		this.renderSizeOption(50)+
		this.renderSizeOption(60)+
		this.renderSizeOption(70)+
		this.renderSizeOption(80)+
		this.renderSizeOption(90)+
		this.renderSizeOption(100)+`
</select>`)

	return `<div class="page">` + strings.Join(result, "") + `</div>`
}

// 判断是否为最后一页
func (this *Page) IsLastPage() bool {
	return this.Current == this.Max
}

func (this *Page) composeURL(page int) string {
	this.Query["page"] = []string{strconv.Itoa(page)}
	return this.Path + "?" + this.Query.Encode()
}

func (this *Page) min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

func (this *Page) max(i, j int) int {
	if i < j {
		return j
	}
	return i
}

func (this *Page) renderSizeOption(size int) string {
	o := `<option value="` + strconv.Itoa(size) + `"`
	if size == this.Size {
		o += ` selected="selected"`
	}
	o += `>` + strconv.Itoa(size) + `条</option>`
	return o
}
