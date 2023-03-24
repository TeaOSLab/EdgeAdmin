// Copyright 2023 Liuxiangchao iwind.liu@gmail.com. All rights reserved. Official site: https://goedge.cn .

package certs

import (
	"bytes"
	"crypto/tls"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/types"
	"io"
	"mime/multipart"
	"strings"
)

type UploadBatchPopupAction struct {
	actionutils.ParentAction
}

func (this *UploadBatchPopupAction) Init() {
	this.Nav("", "", "")
}

func (this *UploadBatchPopupAction) RunGet(params struct{}) {
	this.Data["maxFiles"] = this.maxFiles()

	this.Show()
}

func (this *UploadBatchPopupAction) RunPost(params struct {
	UserId int64

	Must *actions.Must
	CSRF *actionutils.CSRF
}) {
	defer this.CreateLogInfo("批量上传证书")

	var files = this.Request.MultipartForm.File["certFiles"]
	if len(files) == 0 {
		this.Fail("请选择要上传的证书和私钥文件")
		return
	}

	// 限制每次上传的文件数量
	var maxFiles = this.maxFiles()
	if len(files) > maxFiles {
		this.Fail("每次上传最多不能超过" + types.String(maxFiles) + "个文件")
		return
	}

	type certInfo struct {
		filename string
		data     []byte
	}

	var certDataList = []*certInfo{}
	var keyDataList = [][]byte{}

	var failMessages = []string{}
	for _, file := range files {
		func(file *multipart.FileHeader) {
			fp, err := file.Open()
			if err != nil {
				failMessages = append(failMessages, "文件"+file.Filename+"读取失败："+err.Error())
				return
			}

			defer func() {
				_ = fp.Close()
			}()

			data, err := io.ReadAll(fp)
			if err != nil {
				failMessages = append(failMessages, "文件"+file.Filename+"读取失败："+err.Error())
				return
			}

			if bytes.Contains(data, []byte("CERTIFICATE-")) {
				certDataList = append(certDataList, &certInfo{
					filename: file.Filename,
					data:     data,
				})
			} else if bytes.Contains(data, []byte("PRIVATE KEY-")) {
				keyDataList = append(keyDataList, data)
			} else {
				failMessages = append(failMessages, "文件"+file.Filename+"读取失败：文件格式错误，无法识别是证书还是私钥")
				return
			}
		}(file)
	}

	if len(failMessages) > 0 {
		this.Fail("发生了错误：" + strings.Join(failMessages, "；"))
		return
	}

	// 对比证书和私钥数量是否一致
	if len(certDataList) != len(keyDataList) {
		this.Fail("证书文件数量（" + types.String(len(certDataList)) + "）和私钥文件数量（" + types.String(len(keyDataList)) + "）不一致")
		return
	}

	// 自动匹配
	var pairs = [][2][]byte{}        // [] { cert, key }
	var keyIndexMap = map[int]bool{} // 方便下面跳过已匹配的Key
	for _, cert := range certDataList {
		var found = false
		for keyIndex, keyData := range keyDataList {
			if keyIndexMap[keyIndex] {
				continue
			}

			_, err := tls.X509KeyPair(cert.data, keyData)
			if err == nil {
				found = true
				pairs = append(pairs, [2][]byte{cert.data, keyData})
				keyIndexMap[keyIndex] = true
				break
			}
		}
		if !found {
			this.Fail("找不到" + cert.filename + "对应的私钥")
			return
		}
	}

	// 组织 CertConfig
	var pbCerts = []*pb.CreateSSLCertsRequestCert{}
	for _, pair := range pairs {
		certData, keyData := pair[0], pair[1]

		var certConfig = &sslconfigs.SSLCertConfig{
			IsCA:     false,
			CertData: certData,
			KeyData:  keyData,
		}
		err := certConfig.Init(nil)
		if err != nil {
			this.Fail("证书验证失败：" + err.Error())
			return
		}

		var certName = ""
		if len(certConfig.DNSNames) > 0 {
			certName = certConfig.DNSNames[0]
			if len(certConfig.DNSNames) > 1 {
				certName += "等" + types.String(len(certConfig.DNSNames)) + "个域名"
			}
		}

		pbCerts = append(pbCerts, &pb.CreateSSLCertsRequestCert{
			IsOn:        true,
			Name:        certName,
			Description: "",
			ServerName:  "",
			IsCA:        false,
			CertData:    certData,
			KeyData:     keyData,
			TimeBeginAt: certConfig.TimeBeginAt,
			TimeEndAt:   certConfig.TimeEndAt,
			DnsNames:    certConfig.DNSNames,
			CommonNames: certConfig.CommonNames,
		})
	}

	_, err := this.RPC().SSLCertRPC().CreateSSLCerts(this.AdminContext(), &pb.CreateSSLCertsRequest{
		UserId:   params.UserId,
		SSLCerts: pbCerts,
	})
	if err != nil {
		this.ErrorPage(err)
		return
	}

	this.Data["count"] = len(pbCerts)
	this.Success()
}
