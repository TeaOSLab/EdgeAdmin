package webp

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/nodeconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("webp")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	webpResp, err := this.RPC().NodeClusterRPC().FindEnabledNodeClusterWebPPolicy(this.AdminContext(), &pb.FindEnabledNodeClusterWebPPolicyRequest{NodeClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	if len(webpResp.WebpPolicyJSON) == 0 {
		this.Data["webpPolicy"] = nodeconfigs.DefaultWebPImagePolicy
	} else {
		var config = &nodeconfigs.WebPImagePolicy{}
		err = json.Unmarshal(webpResp.WebpPolicyJSON, config)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.Data["webpPolicy"] = config
	}

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId     int64
	IsOn          bool
	Quality       int
	RequireCache  bool
	MinLengthJSON []byte
	MaxLengthJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.ServerWebP_LogUpdateClusterWebPPolicy, params.ClusterId)

	var config = &nodeconfigs.WebPImagePolicy{
		IsOn:         params.IsOn,
		RequireCache: params.RequireCache,
	}

	if params.Quality < 0 {
		params.Quality = 0
	} else if params.Quality > 100 {
		params.Quality = 100
	}
	config.Quality = params.Quality

	if len(params.MinLengthJSON) > 0 {
		var minLength = &shared.SizeCapacity{}
		err := json.Unmarshal(params.MinLengthJSON, minLength)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		config.MinLength = minLength
	}

	if len(params.MaxLengthJSON) > 0 {
		var maxLength = &shared.SizeCapacity{}
		err := json.Unmarshal(params.MaxLengthJSON, maxLength)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		config.MaxLength = maxLength
	}

	err := config.Init()
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
	}

	configJSON, err := json.Marshal(config)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().NodeClusterRPC().UpdateNodeClusterWebPPolicy(this.AdminContext(), &pb.UpdateNodeClusterWebPPolicyRequest{
		NodeClusterId:  params.ClusterId,
		WebPPolicyJSON: configJSON,
	})

	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
