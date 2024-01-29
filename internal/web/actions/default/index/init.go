package index

import (
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Prefix("").
			GetPost("/", new(IndexAction)).
			GetPost("/index/otp", new(OtpAction)).
			GetPost("/initPassword", new(InitPasswordAction)).
			EndAll()
	})
}
