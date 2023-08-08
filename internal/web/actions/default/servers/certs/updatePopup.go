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

type UpdatePopupAction struct {
	actionutils.ParentAction
}

func (this *UpdatePopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UpdatePopupAction) RunGet(params struct {
	CertId int64
}) {
	certConfigResp, err := this.RPC().SSLCertRPC().FindEnabledSSLCertConfig(this.AdminContext(), &pb.FindEnabledSSLCertConfigRequest{SslCertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var certConfigJSON = certConfigResp.SslCertJSON
	if len(certConfigJSON) == 0 {
		this.NotFound("cert", params.CertId)
		return
	}

	var certConfig = &sslconfigs.SSLCertConfig{}
	err = json.Unmarshal(certConfigJSON, certConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	certConfig.CertData = nil // cert & key 不需要在界面上显示
	certConfig.KeyData = nil
	this.Data["certConfig"] = certConfig

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	CertId int64

	TextMode bool

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
	// 创建日志
	defer this.CreateLogInfo(codes.SSLCert_LogUpdateSSLCert, params.CertId)

	// 查询Cert
	certConfigResp, err := this.RPC().SSLCertRPC().FindEnabledSSLCertConfig(this.AdminContext(), &pb.FindEnabledSSLCertConfigRequest{SslCertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	var certConfigJSON = certConfigResp.SslCertJSON
	if len(certConfigJSON) == 0 {
		this.NotFound("cert", params.CertId)
		return
	}

	var certConfig = &sslconfigs.SSLCertConfig{}
	err = json.Unmarshal(certConfigJSON, certConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 校验参数
	params.Must.
		Field("name", params.Name).
		Require("请输入证书说明")

	if params.TextMode {
		if len(params.CertText) > 0 {
			certConfig.CertData = []byte(params.CertText)
		}

		if !params.IsCA {
			if len(params.KeyText) > 0 {
				certConfig.KeyData = []byte(params.KeyText)
			}
		}
	} else {
		if params.CertFile != nil {
			certConfig.CertData, err = params.CertFile.Read()
			if err != nil {
				this.FailField("certFile", "读取证书文件内容错误，请重新上传")
			}
		}

		if !params.IsCA {
			if params.KeyFile != nil {
				certConfig.KeyData, err = params.KeyFile.Read()
				if err != nil {
					this.FailField("keyFile", "读取私钥文件内容错误，请重新上传")
				}
			}
		}
	}

	// 校验
	certConfig.IsCA = params.IsCA
	err = certConfig.Init(context.TODO())
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
	_, err = this.RPC().SSLCertRPC().UpdateSSLCert(this.AdminContext(), &pb.UpdateSSLCertRequest{
		SslCertId:   params.CertId,
		IsOn:        params.IsOn,
		Name:        params.Name,
		Description: params.Description,
		ServerName:  "",
		IsCA:        params.IsCA,
		CertData:    certConfig.CertData,
		KeyData:     certConfig.KeyData,
		TimeBeginAt: certConfig.TimeBeginAt,
		TimeEndAt:   certConfig.TimeEndAt,
		DnsNames:    certConfig.DNSNames,
		CommonNames: certConfig.CommonNames,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Success()
}
