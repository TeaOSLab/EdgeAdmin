package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
)

type CreateSetPopupAction struct {
	actionutils.ParentAction
}

func (this *CreateSetPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreateSetPopupAction) RunGet(params struct {
	HeaderPolicyId int64
}) {
	this.Data["headerPolicyId"] = params.HeaderPolicyId

	this.Show()
}

func (this *CreateSetPopupAction) RunPost(params struct {
	HeaderPolicyId int64
	Name           string
	Value          string

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "设置请求Header，HeaderPolicyId:%d, Name:%s, Value:%s", params.HeaderPolicyId, params.Name, params.Value)

	params.Must.
		Field("name", params.Name).
		Require("请输入Header名称")

	configResp, err := this.RPC().HTTPHeaderPolicyRPC().FindEnabledHTTPHeaderPolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderPolicyConfigRequest{HeaderPolicyId: params.HeaderPolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	policyConfig := &shared.HTTPHeaderPolicy{}
	err = json.Unmarshal(configResp.HeaderPolicyJSON, policyConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 创建Header
	createHeaderResp, err := this.RPC().HTTPHeaderRPC().CreateHTTPHeader(this.AdminContext(), &pb.CreateHTTPHeaderRequest{
		Name:  params.Name,
		Value: params.Value,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	headerId := createHeaderResp.HeaderId

	// 保存
	refs := policyConfig.SetHeaderRefs
	refs = append(refs, &shared.HTTPHeaderRef{
		IsOn:     true,
		HeaderId: headerId,
	})
	refsJSON, err := json.Marshal(refs)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPHeaderPolicyRPC().UpdateHTTPHeaderPolicySettingHeaders(this.AdminContext(), &pb.UpdateHTTPHeaderPolicySettingHeadersRequest{
		HeaderPolicyId: params.HeaderPolicyId,
		HeadersJSON:    refsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
