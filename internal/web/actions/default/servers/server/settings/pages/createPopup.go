package pages

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Status    string
	URL       string `alias:"url"`
	NewStatus int
	Must      *actions.Must
}) {
	// TODO 对状态码进行更多校验

	params.Must.
		Field("status", params.Status).
		Require("请输入响应状态码").
		Field("url", params.URL).
		Require("请输入要显示的URL")

	createResp, err := this.RPC().HTTPPageRPC().CreateHTTPPage(this.AdminContext(), &pb.CreateHTTPPageRequest{
		StatusList: []string{params.Status},
		Url:        params.URL,
		NewStatus:  types.Int32(params.NewStatus),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	pageId := createResp.PageId

	configResp, err := this.RPC().HTTPPageRPC().FindEnabledHTTPPageConfig(this.AdminContext(), &pb.FindEnabledHTTPPageConfigRequest{PageId: pageId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	pageConfig := &serverconfigs.HTTPPageConfig{}
	err = json.Unmarshal(configResp.Config, pageConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["page"] = pageConfig

	this.Success()
}
