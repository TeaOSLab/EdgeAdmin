package pages

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
)

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	PageId int64
}) {
	configResp, err := this.RPC().HTTPPageRPC().FindEnabledHTTPPageConfig(this.AdminContext(), &pb.FindEnabledHTTPPageConfigRequest{PageId: params.PageId})
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
	this.Data["pageConfig"] = pageConfig

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	PageId int64

	Status    string
	URL       string `alias:"url"`
	NewStatus int

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLog(oplogs.LevelInfo, "修改特殊页面 %d", params.PageId)

	params.Must.
		Field("status", params.Status).
		Require("请输入响应状态码").
		Field("url", params.URL).
		Require("请输入要显示的URL")

	_, err := this.RPC().HTTPPageRPC().UpdateHTTPPage(this.AdminContext(), &pb.UpdateHTTPPageRequest{
		PageId:     params.PageId,
		StatusList: []string{params.Status},
		Url:        params.URL,
		NewStatus:  types.Int32(params.NewStatus),
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 返回修改后的配置
	configResp, err := this.RPC().HTTPPageRPC().FindEnabledHTTPPageConfig(this.AdminContext(), &pb.FindEnabledHTTPPageConfigRequest{PageId: params.PageId})
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

	this.Success()
}
