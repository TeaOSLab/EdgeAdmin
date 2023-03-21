package utils

import (
	teaconst "github.com/TeaOSLab/EdgeAdmin/internal/const"
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/iwind/TeaGo/logs"
	"github.com/miekg/dns"
)

var sharedDNSClient *dns.Client
var sharedDNSConfig *dns.ClientConfig

func init() {
	if !teaconst.IsMain {
		return
	}

	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		logs.Println("ERROR: configure dns client failed: " + err.Error())
		return
	}

	sharedDNSConfig = config
	sharedDNSClient = &dns.Client{}
}

// LookupCNAME 获取CNAME
func LookupCNAME(host string) (string, error) {
	var m = new(dns.Msg)

	m.SetQuestion(host+".", dns.TypeCNAME)
	m.RecursionDesired = true

	var lastErr error
	for _, serverAddr := range sharedDNSConfig.Servers {
		r, _, err := sharedDNSClient.Exchange(m, configutils.QuoteIP(serverAddr)+":"+sharedDNSConfig.Port)
		if err != nil {
			lastErr = err
			continue
		}
		if len(r.Answer) == 0 {
			continue
		}

		return r.Answer[0].(*dns.CNAME).Target, nil
	}
	return "", lastErr
}
