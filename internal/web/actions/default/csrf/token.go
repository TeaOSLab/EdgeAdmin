package csrf

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/csrf"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"sync"
	"time"
)

var lastTimestamp = int64(0)
var locker sync.Mutex

type TokenAction struct {
	actionutils.ParentAction
}

func (this *TokenAction) Init() {
	this.Nav("", "", "")
}

func (this *TokenAction) RunGet(params struct {
	Auth *helpers.UserShouldAuth
}) {
	locker.Lock()
	defer locker.Unlock()

	defer func() {
		lastTimestamp = time.Now().Unix()
	}()

	// 没有登录，则限制请求速度
	if params.Auth.AdminId() <= 0 && lastTimestamp > 0 && time.Now().Unix()-lastTimestamp <= 0 {
		this.Fail("请求速度过快，请稍后刷新后重试")
	}

	this.Data["token"] = csrf.Generate()
	this.Success()
}
