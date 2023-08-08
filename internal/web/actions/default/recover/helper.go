package recovers

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/iwind/TeaGo/actions"
)

type Helper struct {
}

func (this *Helper) BeforeAction(actionPtr actions.ActionWrapper) (goNext bool) {
	if !teaconst.IsRecoverMode {
		actionPtr.Object().RedirectURL("/")
		return false
	}
	return true
}
