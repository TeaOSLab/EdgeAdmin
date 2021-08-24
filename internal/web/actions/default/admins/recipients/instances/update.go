package instances

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/monitorconfigs"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction struct {
	actionutils.ParentAction
}

func (this *UpdateAction) Init() {
	this.Nav("", "", "update")
}

func (this *UpdateAction) RunGet(params struct {
	InstanceId int64
}) {
	instanceResp, err := this.RPC().MessageMediaInstanceRPC().FindEnabledMessageMediaInstance(this.AdminContext(), &pb.FindEnabledMessageMediaInstanceRequest{MessageMediaInstanceId: params.InstanceId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	instance := instanceResp.MessageMediaInstance
	if instance == nil || instance.MessageMedia == nil {
		this.NotFound("messageMediaInstance", params.InstanceId)
		return
	}

	mediaParams := maps.Map{}
	if len(instance.ParamsJSON) > 0 {
		err = json.Unmarshal(instance.ParamsJSON, &mediaParams)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	var rateConfig = &monitorconfigs.RateConfig{}
	if len(instance.RateJSON) > 0 {
		err = json.Unmarshal(instance.RateJSON, rateConfig)
		if err != nil {
			this.ErrorPage(err)
			return
		}
	}

	this.Data["instance"] = maps.Map{
		"id":   instance.Id,
		"name": instance.Name,
		"isOn": instance.IsOn,
		"media": maps.Map{
			"type": instance.MessageMedia.Type,
			"name": instance.MessageMedia.Name,
		},
		"description": instance.Description,
		"params":      mediaParams,
		"rate":        rateConfig,
		"hashLife":    instance.HashLife,
	}

	this.Show()
}

func (this *UpdateAction) RunPost(params struct {
	InstanceId int64
	Name       string
	MediaType  string

	EmailSmtp     string
	EmailUsername string
	EmailPassword string
	EmailFrom     string

	WebHookURL          string
	WebHookMethod       string
	WebHookHeaderNames  []string
	WebHookHeaderValues []string
	WebHookContentType  string
	WebHookParamNames   []string
	WebHookParamValues  []string
	WebHookBody         string

	ScriptType      string
	ScriptPath      string
	ScriptLang      string
	ScriptCode      string
	ScriptCwd       string
	ScriptEnvNames  []string
	ScriptEnvValues []string

	DingTalkWebHookURL string

	QyWeixinCorporateId string
	QyWeixinAgentId     string
	QyWeixinAppSecret   string
	QyWeixinTextFormat  string

	QyWeixinRobotWebHookURL string
	QyWeixinRobotTextFormat string

	AliyunSmsSign              string
	AliyunSmsTemplateCode      string
	AliyunSmsTemplateVarNames  []string
	AliyunSmsTemplateVarValues []string
	AliyunSmsAccessKeyId       string
	AliyunSmsAccessKeySecret   string

	TelegramToken string

	GroupIds    string
	Description string
	IsOn        bool

	RateMinutes int32
	RateCount   int32

	HashLife int32

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("修改消息媒介 %d", params.InstanceId)

	params.Must.
		Field("name", params.Name).
		Require("请输入媒介名称").
		Field("mediaType", params.MediaType).
		Require("请选择媒介类型")

	options := maps.Map{}

	switch params.MediaType {
	case "email":
		params.Must.
			Field("emailSmtp", params.EmailSmtp).
			Require("请输入SMTP地址").
			Field("emailUsername", params.EmailUsername).
			Require("请输入邮箱账号").
			Field("emailPassword", params.EmailPassword).
			Require("请输入密码或授权码")

		options["smtp"] = params.EmailSmtp
		options["username"] = params.EmailUsername
		options["password"] = params.EmailPassword
		options["from"] = params.EmailFrom
	case "webHook":
		params.Must.
			Field("webHookURL", params.WebHookURL).
			Require("请输入URL地址").
			Match("(?i)^(http|https)://", "URL地址必须以http或https开头").
			Field("webHookMethod", params.WebHookMethod).
			Require("请选择请求方法")

		options["url"] = params.WebHookURL
		options["method"] = params.WebHookMethod
		options["contentType"] = params.WebHookContentType

		headers := []maps.Map{}
		if len(params.WebHookHeaderNames) > 0 {
			for index, name := range params.WebHookHeaderNames {
				if index < len(params.WebHookHeaderValues) {
					headers = append(headers, maps.Map{
						"name":  name,
						"value": params.WebHookHeaderValues[index],
					})
				}
			}
		}
		options["headers"] = headers

		if params.WebHookContentType == "params" {
			webHookParams := []maps.Map{}
			for index, name := range params.WebHookParamNames {
				if index < len(params.WebHookParamValues) {
					webHookParams = append(webHookParams, maps.Map{
						"name":  name,
						"value": params.WebHookParamValues[index],
					})
				}
			}
			options["params"] = webHookParams
		} else if params.WebHookContentType == "body" {
			options["body"] = params.WebHookBody
		}
	case "script":
		if params.ScriptType == "path" {
			params.Must.
				Field("scriptPath", params.ScriptPath).
				Require("请输入脚本路径")
		} else if params.ScriptType == "code" {
			params.Must.
				Field("scriptCode", params.ScriptCode).
				Require("请输入脚本代码")
		} else {
			params.Must.
				Field("scriptPath", params.ScriptPath).
				Require("请输入脚本路径")
		}

		options["scriptType"] = params.ScriptType
		options["path"] = params.ScriptPath
		options["scriptLang"] = params.ScriptLang
		options["script"] = params.ScriptCode
		options["cwd"] = params.ScriptCwd

		env := []maps.Map{}
		for index, envName := range params.ScriptEnvNames {
			if index < len(params.ScriptEnvValues) {
				env = append(env, maps.Map{
					"name":  envName,
					"value": params.ScriptEnvValues[index],
				})
			}
		}
		options["env"] = env
	case "dingTalk":
		params.Must.
			Field("dingTalkWebHookURL", params.DingTalkWebHookURL).
			Require("请输入Hook地址").
			Match("^https:", "Hook地址必须以https://开头")

		options["webHookURL"] = params.DingTalkWebHookURL
	case "qyWeixin":
		params.Must.
			Field("qyWeixinCorporateId", params.QyWeixinCorporateId).
			Require("请输入企业ID").
			Field("qyWeixinAgentId", params.QyWeixinAgentId).
			Require("请输入应用AgentId").
			Field("qyWeixinSecret", params.QyWeixinAppSecret).
			Require("请输入应用Secret")

		options["corporateId"] = params.QyWeixinCorporateId
		options["agentId"] = params.QyWeixinAgentId
		options["appSecret"] = params.QyWeixinAppSecret
		options["textFormat"] = params.QyWeixinTextFormat
	case "qyWeixinRobot":
		params.Must.
			Field("qyWeixinRobotWebHookURL", params.QyWeixinRobotWebHookURL).
			Require("请输入WebHook地址").
			Match("^https:", "WebHook地址必须以https://开头")

		options["webHookURL"] = params.QyWeixinRobotWebHookURL
		options["textFormat"] = params.QyWeixinRobotTextFormat
	case "aliyunSms":
		params.Must.
			Field("aliyunSmsSign", params.AliyunSmsSign).
			Require("请输入签名名称").
			Field("aliyunSmsTemplateCode", params.AliyunSmsTemplateCode).
			Require("请输入模板CODE").
			Field("aliyunSmsAccessKeyId", params.AliyunSmsAccessKeyId).
			Require("请输入AccessKey ID").
			Field("aliyunSmsAccessKeySecret", params.AliyunSmsAccessKeySecret).
			Require("请输入AccessKey Secret")

		options["sign"] = params.AliyunSmsSign
		options["templateCode"] = params.AliyunSmsTemplateCode
		options["accessKeyId"] = params.AliyunSmsAccessKeyId
		options["accessKeySecret"] = params.AliyunSmsAccessKeySecret

		variables := []maps.Map{}
		for index, name := range params.AliyunSmsTemplateVarNames {
			if index < len(params.AliyunSmsTemplateVarValues) {
				variables = append(variables, maps.Map{
					"name":  name,
					"value": params.AliyunSmsTemplateVarValues[index],
				})
			}
		}
		options["variables"] = variables
	case "telegram":
		params.Must.
			Field("telegramToken", params.TelegramToken).
			Require("请输入机器人Token")
		options["token"] = params.TelegramToken
	default:
		this.Fail("错误的媒介类型")
	}

	optionsJSON, err := json.Marshal(options)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	var rateConfig = &monitorconfigs.RateConfig{
		Minutes: params.RateMinutes,
		Count:   params.RateCount,
	}
	rateJSON, err := json.Marshal(rateConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	_, err = this.RPC().MessageMediaInstanceRPC().UpdateMessageMediaInstance(this.AdminContext(), &pb.UpdateMessageMediaInstanceRequest{
		MessageMediaInstanceId: params.InstanceId,
		Name:                   params.Name,
		MediaType:              params.MediaType,
		ParamsJSON:             optionsJSON,
		Description:            params.Description,
		RateJSON:               rateJSON,
		HashLife:               params.HashLife,
		IsOn:                   params.IsOn,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
