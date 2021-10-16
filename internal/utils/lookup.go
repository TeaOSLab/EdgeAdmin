package utils

import (
	"github.com/TeaOSLab/EdgeCommon/pkg/configutils"
	"github.com/miekg/dns"
)

// LookupCNAME 获取CNAME
func LookupCNAME(host string) (string, error) {
	config, err := dns.ClientConfigFromFile("/etc/resolv.conf")
	if err != nil {
		return "", err
	}

	c := new(dns.Client)
	m := new(dns.Msg)

	m.SetQuestion(host+".", dns.TypeCNAME)
	m.RecursionDesired = true

	var lastErr error
	for _, serverAddr := range config.Servers {
		r, _, err := c.Exchange(m, configutils.QuoteIP(serverAddr)+":"+config.Port)
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
