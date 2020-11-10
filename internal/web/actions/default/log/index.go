package log

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
	"strings"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("log", "log", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().LogRPC().CountLogs(this.AdminContext(), &pb.CountLogRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	logsResp, err := this.RPC().LogRPC().ListLogs(this.AdminContext(), &pb.ListLogsRequest{
		Offset: page.Offset,
		Size:   page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	logMaps := []maps.Map{}
	for _, log := range logsResp.Logs {
		regionName := ""
		log.Ip = "123.123.88.220" // TODO
		regionResp, err := this.RPC().IPLibraryRPC().LookupIPRegion(this.AdminContext(), &pb.LookupIPRegionRequest{Ip: log.Ip})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if regionResp.Region != nil {
			pieces := []string{}
			if len(regionResp.Region.Country) > 0 {
				pieces = append(pieces, regionResp.Region.Country)
			}
			if len(regionResp.Region.Province) > 0 && !lists.ContainsString(pieces, regionResp.Region.Province) {
				pieces = append(pieces, regionResp.Region.Province)
			}
			if len(regionResp.Region.City) > 0 && !lists.ContainsString(pieces, regionResp.Region.City) && !lists.ContainsString(pieces, strings.TrimSuffix(regionResp.Region.Province, "市")) {
				pieces = append(pieces, regionResp.Region.City)
			}
			if len(regionResp.Region.Isp) > 0 && !lists.ContainsString(pieces, regionResp.Region.Isp) {
				pieces = append(pieces, regionResp.Region.Isp)
			}
			regionName = strings.Join(pieces, " ")
		}

		logMaps = append(logMaps, maps.Map{
			"description": log.Description,
			"userName":    log.UserName,
			"createdTime": timeutil.FormatTime("Y-m-d H:i:s", log.CreatedAt),
			"level":       log.Level,
			"type":        log.Type,
			"ip":          log.Ip,
			"region":      regionName,
		})
	}
	this.Data["logs"] = logMaps

	this.Show()
}
