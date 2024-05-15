package cache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/default/nodes/nodeutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/messageconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"net/http"
	"strconv"
)

type TestWriteAction struct {
	actionutils.ParentAction
}

func (this *TestWriteAction) RunPost(params struct {
	ClusterId     int64
	CachePolicyId int64
	Key           string
	Value         string
}) {
	// 记录clusterId
	this.AddCookie(&http.Cookie{
		Name:  "cache_cluster_id",
		Value: strconv.FormatInt(params.ClusterId, 10),
	})

	cachePolicyResp, err := this.RPC().HTTPCachePolicyRPC().FindEnabledHTTPCachePolicyConfig(this.AdminContext(), &pb.FindEnabledHTTPCachePolicyConfigRequest{HttpCachePolicyId: params.CachePolicyId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	cachePolicyJSON := cachePolicyResp.HttpCachePolicyJSON
	if len(cachePolicyJSON) == 0 {
		this.Fail("找不到要操作的缓存策略")
	}

	// 发送命令
	msg := &messageconfigs.WriteCacheMessage{
		CachePolicyJSON: cachePolicyJSON,
		Key:             params.Key,
		Value:           []byte(params.Value),
		LifeSeconds:     3600,
	}
	results, err := nodeutils.SendMessageToCluster(this.AdminContext(), params.ClusterId, messageconfigs.MessageCodeWriteCache, msg, 10, false)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	isAllOk := true
	for _, result := range results {
		if !result.IsOK {
			isAllOk = false
			break
		}
	}

	this.Data["isAllOk"] = isAllOk
	this.Data["results"] = results

	// 创建日志
	defer this.CreateLogInfo(codes.ServerCachePolicy_LogTestWriting, params.CachePolicyId)

	this.Success()
}
