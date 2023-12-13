package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/dns/domains/domainutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/dao"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
}

func (this *IndexAction) RunGet(params struct {
	ServerId   int64
	LocationId int64
}) {
	webConfig, err := dao.SharedHTTPWebDAO.FindWebConfigWithLocationId(this.AdminContext(), params.LocationId)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["webId"] = webConfig.Id
	this.Data["cacheConfig"] = webConfig.Cache

	// 当前集群的缓存策略
	cachePolicy, err := dao.SharedHTTPCachePolicyDAO.FindEnabledHTTPCachePolicyWithServerId(this.AdminContext(), params.ServerId)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if cachePolicy != nil {
		var maxBytes = &shared.SizeCapacity{}
		if !utils.JSONIsNull(cachePolicy.MaxBytesJSON) {
			err = json.Unmarshal(cachePolicy.MaxBytesJSON, maxBytes)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		this.Data["cachePolicy"] = maps.Map{
			"id":       cachePolicy.Id,
			"name":     cachePolicy.Name,
			"isOn":     cachePolicy.IsOn,
			"maxBytes": maxBytes,
		}
	} else {
		this.Data["cachePolicy"] = nil
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	WebId     int64
	CacheJSON []byte

	Must *actions.Must
}) {
	defer this.CreateLogInfo(codes.ServerCache_LogUpdateCacheSettings, params.WebId)

	// 校验配置
	var cacheConfig = &serverconfigs.HTTPCacheConfig{}
	err := json.Unmarshal(params.CacheJSON, cacheConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 检查Key
	if cacheConfig.Key != nil && cacheConfig.Key.IsOn {
		if cacheConfig.Key.Scheme != "http" && cacheConfig.Key.Scheme != "https" {
			this.Fail("缓存主域名协议只能是http或者https")
			return
		}
		if len(cacheConfig.Key.Host) == 0 {
			this.Fail("请输入缓存主域名")
			return
		}
		cacheConfig.Key.Host = strings.ToLower(strings.TrimSuffix(cacheConfig.Key.Host, "/"))
		if !domainutils.ValidateDomainFormat(cacheConfig.Key.Host) {
			this.Fail("请输入正确的缓存主域名")
			return
		}

		// 检查域名所属
		serverIdResp, err := this.RPC().HTTPWebRPC().FindServerIdWithHTTPWebId(this.AdminContext(), &pb.FindServerIdWithHTTPWebIdRequest{HttpWebId: params.WebId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var serverId = serverIdResp.ServerId
		if serverId <= 0 {
			this.Fail("找不到要操作的网站")
			return
		}

		existServerNameResp, err := this.RPC().ServerRPC().CheckServerNameInServer(this.AdminContext(), &pb.CheckServerNameInServerRequest{
			ServerId:   serverId,
			ServerName: cacheConfig.Key.Host,
		})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		if !existServerNameResp.Exists {
			this.Fail("域名 '" + cacheConfig.Key.Host + "' 在当前网站中并未绑定，不能作为缓存主域名")
			return
		}
	}

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
