package issues

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "", "")
}

func (this *IndexAction) RunGet(params struct{}) {
	this.Data["issues"] = []interface{}{}
	this.Show()
}

func (this *IndexAction) RunPost(params struct{}) {
	issuesResp, err := this.RPC().DNSRPC().FindAllDNSIssues(this.AdminContext(), &pb.FindAllDNSIssuesRequest{})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	issueMaps := []maps.Map{}
	for _, issue := range issuesResp.Issues {
		issueMaps = append(issueMaps, maps.Map{
			"target":      issue.Target,
			"targetId":    issue.TargetId,
			"type":        issue.Type,
			"description": issue.Description,
			"params":      issue.Params,
		})
	}
	this.Data["issues"] = issueMaps

	this.Success()
}
