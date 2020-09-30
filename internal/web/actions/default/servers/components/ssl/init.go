package ssl

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/components/componentutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/helpers"
	"github.com/iwind/TeaGo"
)

func init() {
	TeaGo.BeforeStart(func(server *TeaGo.Server) {
		server.
			Helper(helpers.NewUserMustAuth()).
			Helper(NewHelper()).
			Helper(componentutils.NewComponentHelper()).
			Prefix("/servers/components/ssl").
			Get("", new(IndexAction)).
			GetPost("/uploadPopup", new(UploadPopupAction)).
			Post("/delete", new(DeleteAction)).
			GetPost("/updatePopup", new(UpdatePopupAction)).
			Get("/certPopup", new(CertPopupAction)).
			Get("/viewKey", new(ViewKeyAction)).
			Get("/viewCert", new(ViewCertAction)).
			Get("/downloadKey", new(DownloadKeyAction)).
			Get("/downloadCert", new(DownloadCertAction)).
			Get("/downloadZip", new(DownloadZipAction)).
			EndAll()
	})
}
