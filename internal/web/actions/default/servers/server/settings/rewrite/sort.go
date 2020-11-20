package rewrite

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/webutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

type SortAction struct {
	actionutils.ParentAction
}

func (this *SortAction) RunPost(params struct {
	WebId          int64
	RewriteRuleIds []int64
}) {
	defer this.CreateLogInfo("对Web %d 中的重写规则进行排序", params.WebId)

	webConfig, err := webutils.FindWebConfigWithId(this.Parent(), params.WebId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	refsMap := map[int64]*serverconfigs.HTTPRewriteRef{}
	for _, ref := range webConfig.RewriteRefs {
		refsMap[ref.RewriteRuleId] = ref
	}
	newRefs := []*serverconfigs.HTTPRewriteRef{}
	for _, rewriteRuleId := range params.RewriteRuleIds {
		ref, ok := refsMap[rewriteRuleId]
		if ok {
			newRefs = append(newRefs, ref)
		}
	}
	refsJSON, err := json.Marshal(newRefs)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebRewriteRules(this.AdminContext(), &pb.UpdateHTTPWebRewriteRulesRequest{
		WebId:            params.WebId,
		RewriteRulesJSON: refsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
