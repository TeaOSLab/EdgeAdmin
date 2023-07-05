package pages

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
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

		if len(params.Body) > 32*1024 {
			this.FailField("body", "自定义页面内容不能超过32K")
			return
		}
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
	var pageId = createResp.HttpPageId

	configResp, err := this.RPC().HTTPPageRPC().FindEnabledHTTPPageConfig(this.AdminContext(), &pb.FindEnabledHTTPPageConfigRequest{HttpPageId: pageId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var pageConfig = &serverconfigs.HTTPPageConfig{}
	err = json.Unmarshal(configResp.PageJSON, pageConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	err = pageConfig.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
		return
	}

	this.Data["page"] = pageConfig

	// 日志
	defer this.CreateLogInfo(codes.ServerPage_LogCreatePage, pageId)

	this.Success()
}
