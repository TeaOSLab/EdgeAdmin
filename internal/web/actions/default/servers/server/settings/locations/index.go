package locations

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("locations")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["webId"] = webConfig.Id

	if webConfig.Locations != nil {
		this.Data["locations"] = webConfig.Locations
	} else {
		this.Data["locations"] = []interface{}{}
	}

	this.Show()
}
