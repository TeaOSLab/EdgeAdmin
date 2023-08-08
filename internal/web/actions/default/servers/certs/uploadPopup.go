package certs

import (
	"context"
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/langs/codes"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type UploadPopupAction struct {
	actionutils.ParentAction
}

func (this *UploadPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UploadPopupAction) RunGet(params struct {
	ServerId int64
	UserId   int64
}) {
	// 读取服务用户
	if params.ServerId > 0 {
		serverResp, err := this.RPC().ServerRPC().FindEnabledUserServerBasic(this.AdminContext(), &pb.FindEnabledUserServerBasicRequest{ServerId: params.ServerId})
		if err != nil {
			this.ErrorPage(err)
			return
		}
		var server = serverResp.Server
		if server != nil {
			params.UserId = server.UserId
		}
	}
	this.Data["userId"] = params.UserId

	this.Show()
}

func (this *UploadPopupAction) RunPost(params struct {
	UserId int64

	TextMode    bool
	Name        string
	IsCA        bool
	Description string
	IsOn        bool

	CertFile *actions.File
	KeyFile  *actions.File

	CertText string
	KeyText  string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入证书说明")

	var certData []byte
	var keyData []byte

	if params.TextMode {
		if len(params.CertText) == 0 {
			this.FailField("certText", "请输入证书内容")
		}

		if !params.IsCA {
			if len(params.KeyText) == 0 {
				this.FailField("keyText", "请输入私钥内容")
			}
		}

		certData = []byte(params.CertText)
		keyData = []byte(params.KeyText)
	} else {
		if params.CertFile == nil {
			this.FailField("certFile", "请选择要上传的证书文件")
		}
		var err error
		certData, err = params.CertFile.Read()
		if err != nil {
			this.FailField("certFile", "读取证书文件内容错误，请重新上传")
		}

		if !params.IsCA {
			if params.KeyFile == nil {
				this.FailField("keyFile", "请选择要上传的私钥文件")
			} else {
				keyData, err = params.KeyFile.Read()
				if err != nil {
					this.FailField("keyFile", "读取密钥文件内容错误，请重新上传")
				}
			}
		}
	}

	// 校验
	var certConfig = &sslconfigs.SSLCertConfig{
		IsCA:     params.IsCA,
		CertData: certData,
		KeyData:  keyData,
	}
	err := certConfig.Init(context.TODO())
	if err != nil {
		if params.IsCA {
			this.Fail("证书校验错误：" + err.Error())
		} else {
			this.Fail("证书或密钥校验错误：" + err.Error())
		}
	}
	if len(timeutil.Format("Y", certConfig.TimeEnd())) != 4 {
		this.Fail("证书格式错误：无法读取到证书有效期")
	}

	if certConfig.TimeBeginAt < 0 {
		this.Fail("证书校验错误：有效期开始时间过小，不能小于1970年1月1日")
	}
	if certConfig.TimeEndAt < 0 {
		this.Fail("证书校验错误：有效期结束时间过小，不能小于1970年1月1日")
	}

	// 保存
	createResp, err := this.RPC().SSLCertRPC().CreateSSLCert(this.AdminContext(), &pb.CreateSSLCertRequest{
		IsOn:        params.IsOn,
		UserId:      params.UserId,
		Name:        params.Name,
		Description: params.Description,
		ServerName:  "",
		IsCA:        params.IsCA,
		CertData:    certData,
		KeyData:     keyData,
		TimeBeginAt: certConfig.TimeBeginAt,
		TimeEndAt:   certConfig.TimeEndAt,
		DnsNames:    certConfig.DNSNames,
		CommonNames: certConfig.CommonNames,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 查询已创建的证书并返回，方便调用者进行后续处理
	var certId = createResp.SslCertId
	configResp, err := this.RPC().SSLCertRPC().FindEnabledSSLCertConfig(this.AdminContext(), &pb.FindEnabledSSLCertConfigRequest{SslCertId: certId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	certConfig = &sslconfigs.SSLCertConfig{}
	err = json.Unmarshal(configResp.SslCertJSON, certConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	certConfig.CertData = nil // 去掉不必要的数据
	certConfig.KeyData = nil  // 去掉不必要的数据
	this.Data["cert"] = certConfig
	this.Data["certRef"] = &sslconfigs.SSLCertRef{
		IsOn:   true,
		CertId: certId,
	}

	// 创建日志
	defer this.CreateLogInfo(codes.SSLCert_LogUploadSSLCert, certId)

	this.Success()
}
