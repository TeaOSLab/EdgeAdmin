package rewrite

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/server/settings/webutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
)

type DeleteAction struct {
	actionutils.ParentAction
}

func (this *DeleteAction) RunPost(params struct {
	WebId         int64
	RewriteRuleId int64
}) {
	webConfig, err := webutils.FindWebConfigWithId(this.Parent(), params.WebId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	refs := []*serverconfigs.HTTPRewriteRef{}
	for _, ref := range webConfig.RewriteRefs {
		if ref.RewriteRuleId == params.RewriteRuleId {
			continue
		}
		refs = append(refs, ref)
	}

	refsJSON, err := json.Marshal(refs)
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
