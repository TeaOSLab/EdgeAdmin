package waf

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/firewallconfigs"
	"github.com/iwind/TeaGo/actions"
)

type ImportAction struct {
	actionutils.ParentAction
}

func (this *ImportAction) Init() {
	this.Nav("", "", "import")
}

func (this *ImportAction) RunGet(params struct{}) {
	this.Show()
}

func (this *ImportAction) RunPost(params struct {
	FirewallPolicyId int64
	File             *actions.File

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo(codes.WAFPolicy_LogImportWAFPolicy, params.FirewallPolicyId)

	if params.File == nil {
		this.Fail("请上传要导入的文件")
	}
	if params.File.Ext != ".json" {
		this.Fail("规则文件的扩展名只能是.json")
	}

	data, err := params.File.Read()
	if err != nil {
		this.Fail("读取文件时发生错误：" + err.Error())
	}

	config := &firewallconfigs.HTTPFirewallPolicy{}
	err = json.Unmarshal(data, config)
	if err != nil {
		this.Fail("解析文件时发生错误：" + err.Error())
	}

	_, err = this.RPC().HTTPFirewallPolicyRPC().ImportHTTPFirewallPolicy(this.AdminContext(), &pb.ImportHTTPFirewallPolicyRequest{
		HttpFirewallPolicyId:   params.FirewallPolicyId,
		HttpFirewallPolicyJSON: data,
	})
	if err != nil {
		this.Fail("导入失败：" + err.Error())
	}

	this.Success()
}
