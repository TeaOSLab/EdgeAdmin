// Copyright 2022 GoEdge CDN goedge.cdn@gmail.com. All rights reserved. Official site: https://goedge.cn .

package utils

import (
	"github.com/iwind/TeaGo/lists"
	"strings"
)

func FilterNotEmpty(item string) bool {
	return len(item) > 0
}

func MapAddPrefixFunc(prefix string) func(item string) string {
	return func(item string) string {
		if !strings.HasPrefix(item, prefix) {
			return prefix + item
		}
		return item
	}
}

type StringsStream struct {
	s []string
}

func NewStringsStream(s []string) *StringsStream {
	return &StringsStream{s: s}
}

func (this *StringsStream) Map(f ...func(item string) string) *StringsStream {
	for index, item := range this.s {
		for _, f1 := range f {
			item = f1(item)
		}
		this.s[index] = item
	}
	return this
}

func (this *StringsStream) Filter(f ...func(item string) bool) *StringsStream {
	for _, f1 := range f {
		var newStrings = []string{}
		for _, item := range this.s {
			if f1(item) {
				newStrings = append(newStrings, item)
			}
		}
		this.s = newStrings
	}
	return this
}

func (this *StringsStream) Unique() *StringsStream {
	var newStrings = []string{}
	for _, item := range this.s {
		if !lists.ContainsString(newStrings, item) {
			newStrings = append(newStrings, item)
		}
	}
	this.s = newStrings
	return this
}

func (this *StringsStream) Result() []string {
	return this.s
}
