package acme

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "task")
	this.SecondMenu("list")
}

func (this *IndexAction) RunGet(params struct{}) {
	countResp, err := this.RPC().ACMETaskRPC().CountAllEnabledACMETasks(this.AdminContext(), &pb.CountAllEnabledACMETasksRequest{
		AdminId: this.AdminId(),
		UserId:  0,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	count := countResp.Count
	page := this.NewPage(count)
	this.Data["page"] = page.AsHTML()

	tasksResp, err := this.RPC().ACMETaskRPC().ListEnabledACMETasks(this.AdminContext(), &pb.ListEnabledACMETasksRequest{
		AdminId: this.AdminId(),
		UserId:  0,
		Offset:  page.Offset,
		Size:    page.Size,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	taskMaps := []maps.Map{}
	for _, task := range tasksResp.AcmeTasks {
		if task.AcmeUser == nil || task.DnsProvider == nil {
			continue
		}

		var certMap maps.Map = nil
		if task.SslCert != nil {
			certMap = maps.Map{
				"id":   task.SslCert.Id,
				"name": task.SslCert.Name,
			}
		}

		taskMaps = append(taskMaps, maps.Map{
			"id": task.Id,
			"acmeUser": maps.Map{
				"id":    task.AcmeUser.Id,
				"email": task.AcmeUser.Email,
			},
			"dnsProvider": maps.Map{
				"id":   task.DnsProvider.Id,
				"name": task.DnsProvider.Name,
			},
			"dnsDomain": task.DnsDomain,
			"domains":   task.Domains,
			"autoRenew": task.AutoRenew,
			"cert":      certMap,
		})
	}
	this.Data["tasks"] = taskMaps

	this.Show()
}
