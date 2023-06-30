package headers

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/shared"
	"github.com/iwind/TeaGo/actions"
	"strings"
)

type UpdateSetPopupAction struct {
	actionutils.ParentAction
}

func (this *UpdateSetPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdateSetPopupAction) RunGet(params struct {
	HeaderPolicyId int64
	HeaderId       int64
	Type           string
}) {
	this.Data["headerPolicyId"] = params.HeaderPolicyId
	this.Data["headerId"] = params.HeaderId
	this.Data["type"] = params.Type

	headerResp, err := this.RPC().HTTPHeaderRPC().FindEnabledHTTPHeaderConfig(this.AdminContext(), &pb.FindEnabledHTTPHeaderConfigRequest{HeaderId: params.HeaderId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	headerConfig := &shared.HTTPHeaderConfig{}
	err = json.Unmarshal(headerResp.HeaderJSON, headerConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["headerConfig"] = headerConfig

	this.Show()
}

func (this *UpdateSetPopupAction) RunPost(params struct {
	HeaderId int64
	Name     string
	Value    string

	StatusListJSON    []byte
	MethodsJSON       []byte
	DomainsJSON       []byte
	ShouldAppend      bool
	DisableRedirect   bool
	ShouldReplace     bool
	ReplaceValuesJSON []byte

	Must *actions.Must
}) {
	// 日志
	defer this.CreateLogInfo(codes.ServerHTTPHeader_LogUpdateSettingHeader, params.HeaderId, params.Name, params.Value)

	params.Name = strings.TrimSuffix(params.Name, ":")

	params.Must.
		Field("name", params.Name).
		Require("请输入Header名称")

	// status list
	var statusList = []int32{}
	if len(params.StatusListJSON) > 0 {
		err := json.Unmarshal(params.StatusListJSON, &statusList)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// methods
	var methods = []string{}
	if len(params.MethodsJSON) > 0 {
		err := json.Unmarshal(params.MethodsJSON, &methods)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// domains
	var domains = []string{}
	if len(params.DomainsJSON) > 0 {
		err := json.Unmarshal(params.DomainsJSON, &domains)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	// replace values
	var replaceValues = []*shared.HTTPHeaderReplaceValue{}
	if len(params.ReplaceValuesJSON) > 0 {
		err := json.Unmarshal(params.ReplaceValuesJSON, &replaceValues)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	_, err := this.RPC().HTTPHeaderRPC().UpdateHTTPHeader(this.AdminContext(), &pb.UpdateHTTPHeaderRequest{
		HeaderId:          params.HeaderId,
		Name:              params.Name,
		Value:             params.Value,
		Status:            statusList,
		Methods:           methods,
		Domains:           domains,
		ShouldAppend:      params.ShouldAppend,
		DisableRedirect:   params.DisableRedirect,
		ShouldReplace:     params.ShouldReplace,
		ReplaceValuesJSON: params.ReplaceValuesJSON,
	})
	if err != nil {
		this.ErrorPage(err)
	}

	this.Success()
}
