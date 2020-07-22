package index

import (
	"github.com/iwind/TeaGo/actions"
	stringutil "github.com/iwind/TeaGo/utils/string"
)

type IndexAction actions.Action

// 首页（登录页）

var TokenSalt = stringutil.Rand(32)

func (this *IndexAction) RunGet(params struct {
	From string
}) {
	this.WriteString("Hello, i am index")
}
