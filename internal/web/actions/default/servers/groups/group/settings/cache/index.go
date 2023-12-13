package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/servers/groups/group/servergrouputils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "index")
	this.SecondMenu("cache")
}

func (this *IndexAction) RunGet(params struct {
	GroupId int64
}) {
	_, err := servergrouputils.InitGroup(this.Parent(), params.GroupId, "cache")
	if err != nil {
		this.ErrorPage(err)
		return
	}

	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithServerGroupId(this.AdminContext(), params.GroupId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["cacheConfig"] = webConfig.Cache

	this.Data["cachePolicy"] = nil

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId     int64
	CacheJSON []byte

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLogInfo(codes.ServerCache_LogUpdateCacheSettings, params.WebId)

	// 校验配置
	var cacheConfig = &serverconfigs.HTTPCacheConfig{}
	err := json.Unmarshal(params.CacheJSON, cacheConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 分组不支持主域名
	cacheConfig.Key = nil

	err = cacheConfig.Init()
	if err != nil {
		this.Fail("检查配置失败：" + err.Error())
	}

	// 去除不必要的部分
	for _, cacheRef := range cacheConfig.CacheRefs {
		cacheRef.CachePolicy = nil
	}

	cacheJSON, err := json.Marshal(cacheConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().HTTPWebRPC().UpdateHTTPWebCache(this.AdminContext(), &pb.UpdateHTTPWebCacheRequest{
		HttpWebId: params.WebId,
		CacheJSON: cacheJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
