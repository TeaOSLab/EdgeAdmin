package pages

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
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
	this.Data["bodyTypes"] = shared.FindAllBodyTypes()

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Status   string
	BodyType string

	URL  string `alias:"url"`
	Body string

	NewStatus int
	Must      *actions.Must
}) {
	// TODO 对状态码进行更多校验

	params.Must.
		Field("status", params.Status).
		Require("请输入响应状态码")

	switch params.BodyType {
	case shared.BodyTypeURL:
		params.Must.
			Field("url", params.URL).
			Require("请输入要显示的URL")
	case shared.BodyTypeHTML:
		params.Must.
			Field("body", params.Body).
			Require("请输入要显示的HTML内容")
	}

	createResp, err := this.RPC().HTTPPageRPC().CreateHTTPPage(this.AdminContext(), &pb.CreateHTTPPageRequest{
		StatusList: []string{params.Status},
		BodyType:   params.BodyType,
		Url:        params.URL,
		Body:       params.Body,
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
	err = json.Unmarshal(configResp.PageJSON, pageConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["page"] = pageConfig

	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "创建特殊页面 %d", pageId)

	this.Success()
}
