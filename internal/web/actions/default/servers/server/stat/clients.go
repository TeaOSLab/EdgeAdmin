package stat

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type ClientsAction struct {
	actionutils.ParentAction
}

func (this *ClientsAction) Init() {
	this.Nav("", "stat", "")
	this.SecondMenu("client")
}

func (this *ClientsAction) RunGet(params struct {
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
	systemMaps := []maps.Map{}
	browserMaps := []maps.Map{}

	if statIsOn {
		{
			resp, err := this.RPC().ServerClientSystemMonthlyStatRPC().FindTopServerClientSystemMonthlyStats(this.AdminContext(), &pb.FindTopServerClientSystemMonthlyStatsRequest{
				ServerId: params.ServerId,
				Month:    month,
				Offset:   0,
				Size:     10,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			for _, stat := range resp.Stats {
				systemMaps = append(systemMaps, maps.Map{
					"count": stat.Count,
					"system": maps.Map{
						"id":   stat.ClientSystem.Id,
						"name": stat.ClientSystem.Name + " " + stat.Version,
					},
				})
			}
		}

		{
			resp, err := this.RPC().ServerClientBrowserMonthlyStatRPC().FindTopServerClientBrowserMonthlyStats(this.AdminContext(), &pb.FindTopServerClientBrowserMonthlyStatsRequest{
				ServerId: params.ServerId,
				Month:    month,
				Offset:   0,
				Size:     10,
			})
			if err != nil {
				this.ErrorPage(err)
				return
			}
			for _, stat := range resp.Stats {
				browserMaps = append(browserMaps, maps.Map{
					"count": stat.Count,
					"browser": maps.Map{
						"id":   stat.ClientBrowser.Id,
						"name": stat.ClientBrowser.Name + " " + stat.Version,
					},
				})
			}
		}
	}
	this.Data["systemStats"] = systemMaps
	this.Data["browserStats"] = browserMaps

	this.Show()
}
