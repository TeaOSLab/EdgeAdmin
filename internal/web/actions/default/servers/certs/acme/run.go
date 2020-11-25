package acme

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"time"
)

type RunAction struct {
	actionutils.ParentAction
}

func (this *RunAction) RunPost(params struct{}) {
	time.Sleep(5 * time.Second) // TODO
	this.Success()
}
