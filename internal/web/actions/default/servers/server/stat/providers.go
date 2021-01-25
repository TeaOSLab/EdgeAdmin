package stat

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type ProvidersAction struct {
	actionutils.ParentAction
}

func (this *ProvidersAction) Init() {
	this.Nav("", "stat", "")
	this.SecondMenu("provider")
}

func (this *ProvidersAction) RunGet(params struct {
	ServerId int64
	Month    string
}) {
	month := params.Month
	if len(month) != 6 {
		month = timeutil.Format("Ym")
	}
	this.Data["month"] = month

	serverTypeResp, err := this.RPC().ServerRPC().FindEnabledServerType(this.AdminContext(), &pb.FindEnabledServerTypeRequest{ServerId: params.ServerId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	serverType := serverTypeResp.Type

	statIsOn := false

	// 是否已开启
	if serverconfigs.IsHTTPServerType(serverType) {
		webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if webConfig != nil && webConfig.StatRef != nil {
			statIsOn = webConfig.StatRef.IsOn
		}
	} else {
		this.WriteString("此类型服务暂不支持统计")
		return
	}

	this.Data["statIsOn"] = statIsOn

	// 统计数据
	providerMaps := []maps.Map{}

	if statIsOn {
		{
			resp, err := this.RPC().ServerRegionProviderMonthlyStatRPC().FindTopServerRegionProviderMonthlyStats(this.AdminContext(), &pb.FindTopServerRegionProviderMonthlyStatsRequest{
				Month:    month,
				ServerId: params.ServerId,
				Offset:   0,
				Size:     10,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			for _, stat := range resp.Stats {
				providerMaps = append(providerMaps, maps.Map{
					"count": stat.Count,
					"provider": maps.Map{
						"id":   stat.RegionProvider.Id,
						"name": stat.RegionProvider.Name,
					},
				})
			}
		}
	}
	this.Data["providerStats"] = providerMaps

	this.Show()
}
