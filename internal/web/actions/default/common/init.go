package common

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(new(helpers.UserMustAuth)).
			Prefix("/common").
			Get("/changedClusters", new(ChangedClustersAction)).
			Post("/syncClusters", new(SyncClustersAction)).
			EndAll()
	})
}
