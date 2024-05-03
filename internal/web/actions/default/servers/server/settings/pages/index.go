package pages

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	"regexp"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("pages")
}

func (this *IndexAction) RunGet(params struct {
	ServerId int64
}) {
	// 只有HTTP服务才支持
	if this.FilterHTTPFamily() {
		return
	}

	// 分组设置
	groupResp, err := this.RPC().ServerGroupRPC().FindEnabledServerGroupConfigInfo(this.AdminContext(), &pb.FindEnabledServerGroupConfigInfoRequest{
		ServerId: params.ServerId,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["hasGroupConfig"] = groupResp.HasPagesConfig
	this.Data["groupSettingURL"] = "/servers/groups/group/settings/pages?groupId=" + types.String(groupResp.ServerGroupId)

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["enableGlobalPages"] = webConfig.EnableGlobalPages
	this.Data["pages"] = webConfig.Pages
	this.Data["shutdownConfig"] = webConfig.Shutdown

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId             int64
	PagesJSON         []byte
	ShutdownJSON      []byte
	EnableGlobalPages bool
	Must              *actions.Must
}) {
	// 日志
	defer this.CreateLogInfo(codes.ServerPage_LogUpdatePages, params.WebId)

	// 检查配置
	var urlReg = regexp.MustCompile(`^(?i)(http|https)://`)

	// validate pages
	if len(params.PagesJSON) > 0 {
		var pages = []*serverconfigs.HTTPPageConfig{}
		err := json.Unmarshal(params.PagesJSON, &pages)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		for _, page := range pages {
			err = page.Init()
			if err != nil {
				this.Fail("配置校验失败：" + err.Error())
				return
			}

			// check url
			if page.BodyType == serverconfigs.HTTPPageBodyTypeURL && !urlReg.MatchString(page.URL) {
				this.Fail("自定义页面中 '" + page.URL + "' 不是一个正确的URL，请进行修改")
				return
			}
			if page.BodyType == serverconfigs.HTTPPageBodyTypeRedirectURL && !urlReg.MatchString(page.URL) {
				this.Fail("自定义页面中 '" + page.URL + "' 不是一个正确的URL，请进行修改")
				return
			}
		}
	}

	// validate shutdown page
	if len(params.ShutdownJSON) > 0 {
		var shutdownConfig = &serverconfigs.HTTPShutdownConfig{}
		err := json.Unmarshal(params.ShutdownJSON, shutdownConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}

		err = shutdownConfig.Init()
		if err != nil {
			this.Fail("配置校验失败：" + err.Error())
			return
		}

		if shutdownConfig.BodyType == serverconfigs.HTTPPageBodyTypeURL {
			if len(shutdownConfig.URL) > 512 {
				this.Fail("临时关闭页面中URL过长，不能超过512字节")
				return
			}

			if shutdownConfig.IsOn /** 只有启用的时候才校验 **/ && !urlReg.MatchString(shutdownConfig.URL) {
				this.Fail("临时关闭页面中 '" + shutdownConfig.URL + "' 不是一个正确的URL，请进行修改")
				return
			}
		} else if shutdownConfig.BodyType == serverconfigs.HTTPPageBodyTypeRedirectURL {
			if len(shutdownConfig.URL) > 512 {
				this.Fail("临时关闭页面中URL过长，不能超过512字节")
				return
			}

			if shutdownConfig.IsOn /** 只有启用的时候才校验 **/ && !urlReg.MatchString(shutdownConfig.URL) {
				this.Fail("临时关闭页面中 '" + shutdownConfig.URL + "' 不是一个正确的URL，请进行修改")
				return
			}
		} else if shutdownConfig.Body == serverconfigs.HTTPPageBodyTypeHTML {
			if len(shutdownConfig.Body) > 32*1024 {
				this.Fail("临时关闭页面中HTML内容长度不能超过32K")
				return
			}
		}
	}

	_, err := this.RPC().HTTPWebRPC().UpdateHTTPWebGlobalPagesEnabled(this.AdminContext(), &pb.UpdateHTTPWebGlobalPagesEnabledRequest{
		HttpWebId: params.WebId,
		IsEnabled: params.EnableGlobalPages,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebPages(this.AdminContext(), &pb.UpdateHTTPWebPagesRequest{
		HttpWebId: params.WebId,
		PagesJSON: params.PagesJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebShutdown(this.AdminContext(), &pb.UpdateHTTPWebShutdownRequest{
		HttpWebId:    params.WebId,
		ShutdownJSON: params.ShutdownJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
