package components

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs"
	"github.com/iwind/TeaGo/actions"
)

const (
	SettingCodeGlobalConfig = "globalConfig"
)

type IndexAction struct {
	actionutils.ParentAction
}

func (this *IndexAction) Init() {
	this.Nav("", "component", "index")
	this.SecondMenu("global")
}

func (this *IndexAction) RunGet(params struct{}) {
	valueJSONResp, err := this.RPC().SysSettingRPC().ReadSysSetting(this.AdminContext(), &pb.ReadSysSettingRequest{Code: SettingCodeGlobalConfig})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	valueJSON := valueJSONResp.ValueJSON
	globalConfig := &serverconfigs.GlobalConfig{}
	if len(valueJSON) > 0 {
		err = json.Unmarshal(valueJSON, globalConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}
	this.Data["globalConfig"] = globalConfig

	this.Show()
}

func (this *IndexAction) RunPost(params struct {
	GlobalConfigJSON []byte
	Must             *actions.Must
}) {
	if len(params.GlobalConfigJSON) == 0 {
		this.Fail("错误的配置信息，请刷新当前页面后重试")
	}

	globalConfig := &serverconfigs.GlobalConfig{}
	err := json.Unmarshal(params.GlobalConfigJSON, globalConfig)
	if err != nil {
		this.Fail("配置校验失败：" + err.Error())
	}

	// 修改配置
	_, err = this.RPC().SysSettingRPC().UpdateSysSetting(this.AdminContext(), &pb.UpdateSysSettingRequest{
		Code:      SettingCodeGlobalConfig,
		ValueJSON: params.GlobalConfigJSON,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
