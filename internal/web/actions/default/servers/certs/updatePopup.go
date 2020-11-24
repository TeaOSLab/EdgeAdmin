package certs

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/oplogs"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
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
	certConfigResp, err := this.RPC().SSLCertRPC().FindEnabledSSLCertConfig(this.AdminContext(), &pb.FindEnabledSSLCertConfigRequest{CertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	certConfigJSON := certConfigResp.CertJSON
	if len(certConfigJSON) == 0 {
		this.NotFound("cert", params.CertId)
		return
	}

	certConfig := &sslconfigs.SSLCertConfig{}
	err = json.Unmarshal(certConfigJSON, certConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}
	this.Data["certConfig"] = certConfig

	this.Show()
}

func (this *UpdatePopupAction) RunPost(params struct {
	CertId int64

	Name        string
	IsCA        bool
	Description string
	IsOn        bool

	CertFile *actions.File
	KeyFile  *actions.File

	Must *actions.Must
}) {
	// 创建日志
	defer this.CreateLog(oplogs.LevelInfo, "修改SSL证书 %d", params.CertId)

	// 查询Cert
	certConfigResp, err := this.RPC().SSLCertRPC().FindEnabledSSLCertConfig(this.AdminContext(), &pb.FindEnabledSSLCertConfigRequest{CertId: params.CertId})
	if err != nil {
		this.ErrorPage(err)
		return
	}
	certConfigJSON := certConfigResp.CertJSON
	if len(certConfigJSON) == 0 {
		this.NotFound("cert", params.CertId)
		return
	}

	certConfig := &sslconfigs.SSLCertConfig{}
	err = json.Unmarshal(certConfigJSON, certConfig)
	if err != nil {
		this.ErrorPage(err)
		return
	}

	// 校验参数
	params.Must.
		Field("name", params.Name).
		Require("请输入证书说明")

	if params.CertFile != nil {
		certConfig.CertData, err = params.CertFile.Read()
		if err != nil {
			this.Fail("读取证书文件内容错误，请重新上传")
		}
	}

	if !params.IsCA {
		if params.KeyFile != nil {
			certConfig.KeyData, err = params.KeyFile.Read()
			if err != nil {
				this.Fail("读取密钥文件内容错误，请重新上传")
			}
		}
	}

	// 校验
	certConfig.IsCA = params.IsCA
	err = certConfig.Init()
	if err != nil {
		if params.IsCA {
			this.Fail("证书校验错误：" + err.Error())
		} else {
			this.Fail("证书或密钥校验错误：" + err.Error())
		}
	}

	// 保存
	_, err = this.RPC().SSLCertRPC().UpdateSSLCert(this.AdminContext(), &pb.UpdateSSLCertRequest{
		CertId:      params.CertId,
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
