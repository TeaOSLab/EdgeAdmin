package actionutils

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/csrf"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type CSRF struct {
}

func (this *CSRF) BeforeAction(actionPtr actions.ActionWrapper, paramName string) (goNext bool) {
	action := actionPtr.Object()
	token := action.ParamString("csrfToken")
	if !csrf.Validate(token) {
		action.ResponseWriter.WriteHeader(http.StatusForbidden)
		action.WriteString("表单已失效，请刷新页面后重试(001)")
		return
	}

	return true
}
