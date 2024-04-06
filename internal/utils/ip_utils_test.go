package utils

import (
	"net"
	"testing"
)

func TestExtractIP(t *testing.T) {
	t.Log(ExtractIP("192.168.1.100"))
}

func TestExtractIP_CIDR(t *testing.T) {
	t.Log(ExtractIP("192.168.2.100/24"))
}

func TestExtractIP_Range(t *testing.T) {
	t.Log(ExtractIP("192.168.2.100 - 192.168.4.2"))
}

func TestNextIP(t *testing.T) {
	for _, ip := range []string{"192.168.1.1", "0.0.0.0", "255.255.255.255", "192.168.2.255", "192.168.255.255"} {
		t.Log(ip+":", NextIP(net.ParseIP(ip).To4()))
	}
}

func TestNextIP_Copy(t *testing.T) {
	var ip = net.ParseIP("192.168.1.100")
	var nextIP = NextIP(ip)
	t.Log(ip, nextIP)
}
