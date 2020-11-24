package certs

import (
	"encoding/json"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/serverconfigs/sslconfigs"
)

// 所有相关数据
type DatajsAction struct {
	actionutils.ParentAction
}

func (this *DatajsAction) Init() {
}

func (this *DatajsAction) RunGet(params struct{}) {
	this.AddHeader("Content-Type", "text/javascript; charset=utf-8")

	{
		cipherSuitesJSON, err := json.Marshal(sslconfigs.AllTLSCipherSuites)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.WriteString("window.SSL_ALL_CIPHER_SUITES = " + string(cipherSuitesJSON) + ";\n")
	}
	{
		modernCipherSuitesJSON, err := json.Marshal(sslconfigs.TLSModernCipherSuites)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.WriteString("window.SSL_MODERN_CIPHER_SUITES = " + string(modernCipherSuitesJSON) + ";\n")
	}
	{
		intermediateCipherSuitesJSON, err := json.Marshal(sslconfigs.TLSIntermediateCipherSuites)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.WriteString("window.SSL_INTERMEDIATE_CIPHER_SUITES = " + string(intermediateCipherSuitesJSON) + ";\n")
	}
	{
		sslVersionsJSON, err := json.Marshal(sslconfigs.AllTlsVersions)
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.WriteString("window.SSL_ALL_VERSIONS = " + string(sslVersionsJSON) + ";\n")
	}
	{
		clientAuthTypesJSON, err := json.Marshal(sslconfigs.AllSSLClientAuthTypes())
		if err != nil {
			this.ErrorPage(err)
			return
		}
		this.WriteString("window.SSL_ALL_CLIENT_AUTH_TYPES = " + string(clientAuthTypesJSON) + ";\n")
	}
}
