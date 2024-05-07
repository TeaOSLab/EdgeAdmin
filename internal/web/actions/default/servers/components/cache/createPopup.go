package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Data["types"] = serverconfigs.AllCachePolicyStorageTypes

	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name string
	Type string

	// file
	FileDir                string
	FileMemoryCapacityJSON []byte
	FileOpenFileCacheMax   int
	FileEnableSendfile     bool
	FileMinFreeSizeJSON    []byte

	CapacityJSON                   []byte
	MaxSizeJSON                    []byte
	FetchTimeoutJSON               []byte
	SyncCompressionCache           bool
	EnableMMAP                     bool
	EnableIncompletePartialContent bool

	Description string
	IsOn        bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称")

	// 校验选项
	var options any
	switch params.Type {
	case serverconfigs.CachePolicyStorageFile:
		params.Must.
			Field("fileDir", params.FileDir).
			Require("请输入缓存目录")

		var memoryCapacity = &shared.SizeCapacity{}
		if len(params.FileMemoryCapacityJSON) > 0 {
			err := json.Unmarshal(params.FileMemoryCapacityJSON, memoryCapacity)
			if err != nil {
				this.ErrorPage(err)
				return
			}
		}

		var openFileCacheConfig *serverconfigs.OpenFileCacheConfig
		if params.FileOpenFileCacheMax > 0 {
			openFileCacheConfig = &serverconfigs.OpenFileCacheConfig{
				IsOn: true,
				Max:  params.FileOpenFileCacheMax,
			}
		}

		var minFreeSize = &shared.SizeCapacity{}
		var minFreeSizeJSON = params.FileMinFreeSizeJSON
		if !utils.JSONIsNull(minFreeSizeJSON) {
			_, err := utils.JSONDecodeConfig(minFreeSizeJSON, minFreeSize)
			if err != nil {
				this.ErrorPage(err)
				return
			}
			if minFreeSize.Count < 0 {
				minFreeSize.Count = 0
			}
		}

		options = &serverconfigs.HTTPFileCacheStorage{
			Dir: params.FileDir,
			MemoryPolicy: &serverconfigs.HTTPCachePolicy{
				Capacity: memoryCapacity,
			},
			OpenFileCache:                  openFileCacheConfig,
			EnableSendfile:                 params.FileEnableSendfile,
			MinFreeSize:                    minFreeSize,
			EnableMMAP:                     params.EnableMMAP,
			EnableIncompletePartialContent: params.EnableIncompletePartialContent,
		}
	case serverconfigs.CachePolicyStorageMemory:
		options = &serverconfigs.HTTPMemoryCacheStorage{}
	default:
		this.Fail("请选择正确的缓存类型")
	}

	optionsJSON, err := json.Marshal(options)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	createResp, err := this.RPC().HTTPCachePolicyRPC().CreateHTTPCachePolicy(this.AdminContext(), &pb.CreateHTTPCachePolicyRequest{
		IsOn:                 params.IsOn,
		Name:                 params.Name,
		Description:          params.Description,
		CapacityJSON:         params.CapacityJSON,
		MaxSizeJSON:          params.MaxSizeJSON,
		FetchTimeoutJSON:     params.FetchTimeoutJSON,
		Type:                 params.Type,
		OptionsJSON:          optionsJSON,
		SyncCompressionCache: params.SyncCompressionCache,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 返回数据
	this.Data["cachePolicy"] = maps.Map{
		"id":   createResp.HttpCachePolicyId,
		"name": params.Name,
	}

	// 创建日志
	defer this.CreateLogInfo(codes.ServerCachePolicy_LogCreateCachePolicy, createResp.HttpCachePolicyId)

	this.Success()
}
