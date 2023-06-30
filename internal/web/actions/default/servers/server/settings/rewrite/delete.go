package rewrite

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
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
	defer this.CreateLogInfo(codes.HTTPRewriteRule_LogDeleteRewriteRule, params.WebId, params.RewriteRuleId)

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithId(this.AdminContext(), params.WebId)
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
		HttpWebId:        params.WebId,
		RewriteRulesJSON: refsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
