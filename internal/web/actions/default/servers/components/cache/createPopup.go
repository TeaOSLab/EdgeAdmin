package cache

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

type CreatePopupAction struct {
	actionutils.ParentAction
}

func (this *CreatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *CreatePopupAction) RunGet(params struct{}) {
	this.Data["types"] = serverconfigs.AllCachePolicyTypes
	this.Show()
}

func (this *CreatePopupAction) RunPost(params struct {
	Name string
	Type string

	// file
	FileDir string

	CapacityJSON []byte
	MaxSizeJSON  []byte
	MaxKeys      int64

	Description string
	IsOn        bool

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称")

	// 校验选项
	var options interface{}
	switch params.Type {
	case serverconfigs.CachePolicyTypeFile:
		params.Must.
			Field("fileDir", params.FileDir).
			Require("请输入缓存目录")
		options = &serverconfigs.HTTPFileCacheConfig{
			Dir: params.FileDir,
		}
	case serverconfigs.CachePolicyTypeMemory:
		options = &serverconfigs.HTTPMemoryCacheConfig{
		}
	default:
		this.Fail("请选择正确的缓存类型")
	}

	optionsJSON, err := json.Marshal(options)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	_, err = this.RPC().HTTPCachePolicyRPC().CreateHTTPCachePolicy(this.AdminContext(), &pb.CreateHTTPCachePolicyRequest{
		IsOn:         params.IsOn,
		Name:         params.Name,
		Description:  params.Description,
		CapacityJSON: params.CapacityJSON,
		MaxKeys:      params.MaxKeys,
		MaxSizeJSON:  params.MaxSizeJSON,
		Type:         params.Type,
		OptionsJSON:  optionsJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
