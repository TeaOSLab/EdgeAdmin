package locations

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	"strings"
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
	// 只有HTTP服务才支持
	if this.FilterHTTPFamily() {
		return
	}

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["webId"] = webConfig.Id

	var locationMaps = []maps.Map{}
	if webConfig.Locations != nil {
		for _, location := range webConfig.Locations {
			err := location.ExtractPattern()
			if err != nil {
				continue
			}
			jsonData, err := json.Marshal(location)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			m := maps.Map{}
			err = json.Unmarshal(jsonData, &m)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			var pieces = strings.Split(location.Pattern, " ")
			if len(pieces) == 2 {
				m["pattern"] = pieces[1]
				m["patternTypeName"] = serverconfigs.FindLocationPatternTypeName(location.PatternType())
			} else {
				m["pattern"] = location.Pattern
				m["patternTypeName"] = serverconfigs.FindLocationPatternTypeName(serverconfigs.HTTPLocationPatternTypePrefix)
			}
			locationMaps = append(locationMaps, m)
		}
	}
	this.Data["locations"] = locationMaps

	this.Show()
}
