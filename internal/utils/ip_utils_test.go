package utils

import (
	"testing"
)

func TestIP2Long(t *testing.T) {
	for _, ip := range []string{
		"0.0.0.1",
		"0.0.0.2",
		"127.0.0.1",
		"192.0.0.2",
		"255.255.255.255",
		"2001:db8:0:1::101",
		"2001:db8:0:1::102",
		"2406:8c00:0:3409:133:18:203:0",
		"2406:8c00:0:3409:133:18:203:158",
		"2406:8c00:0:3409:133:18:203:159",
		"2406:8c00:0:3409:133:18:203:160",
	} {
		t.Log(ip, " -> ", IP2Long(ip))
	}
}

func TestIsIPv4(t *testing.T) {
	type testIP struct {
		ip string
		ok bool
	}
	for _, item := range []testIP{
		{
			ip: "1",
			ok: false,
		},
		{
			ip: "192.168.0.1",
			ok: true,
		},
		{
			ip: "1.1.0.1",
			ok: true,
		},
		{
			ip: "255.255.255.255",
			ok: true,
		},
		{
			ip: "192.168.0.1233",
			ok: false,
		},
	} {
		if IsIPv4(item.ip) != item.ok {
			t.Fatal(item.ip, "should be", item.ok)
		}
	}
}

func TestIsIPv6(t *testing.T) {
	type testIP struct {
		ip string
		ok bool
	}
	for _, item := range []testIP{
		{
			ip: "1",
			ok: false,
		},
		{
			ip: "2406:8c00:0:3409:133:18:203:158",
			ok: true,
		},
		{
			ip: "2406::8c00:0:3409:133:18:203:158",
			ok: false,
		},
		{
			ip: "2001:db8:0:1::101",
			ok: true,
		},
	} {
		if IsIPv6(item.ip) != item.ok {
			t.Fatal(item.ip, "should be", item.ok)
		}
	}
}
