package utils

import (
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
	r, _, err := c.Exchange(m, config.Servers[0]+":"+config.Port)
	if err != nil {
		return "", err
	}
	if len(r.Answer) == 0 {
		return "", nil
	}

	return r.Answer[0].(*dns.CNAME).Target, nil
}
