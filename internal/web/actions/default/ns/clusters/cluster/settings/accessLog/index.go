// Copyright 2021 Liuxiangchao iwind.liu@gmail.com. All rights reserved.

package accessLog

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/dnsconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "setting", "")
	this.SecondMenu("accessLog")
}

func (this *IndexAction) RunGet(params struct {
	ClusterId int64
}) {
	accessLogResp, err := this.RPC().NSClusterRPC().FindNSClusterAccessLog(this.AdminContext(), &pb.FindNSClusterAccessLogRequest{NsClusterId: params.ClusterId})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	accessLogRef := &dnsconfigs.NSAccessLogRef{}
	if len(accessLogResp.AccessLogJSON) > 0 {
		err = json.Unmarshal(accessLogResp.AccessLogJSON, accessLogRef)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["accessLogRef"] = accessLogRef

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	ClusterId     int64
	AccessLogJSON []byte

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改域名服务集群 %d 访问日志配置", params.ClusterId)

	ref := &dnsconfigs.NSAccessLogRef{}
	err := json.Unmarshal(params.AccessLogJSON, ref)
	if err != nil {
		this.Fail("数据格式错误：" + err.Error())
	}
	err = ref.Init()
	if err != nil {
		this.Fail("数据格式错误：" + err.Error())
	}

	_, err = this.RPC().NSClusterRPC().UpdateNSClusterAccessLog(this.AdminContext(), &pb.UpdateNSClusterAccessLogRequest{
		NsClusterId:   params.ClusterId,
		AccessLogJSON: params.AccessLogJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
