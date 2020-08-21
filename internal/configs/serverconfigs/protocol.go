package serverconfigs

type Protocol = string

const (
	ProtocolHTTP  Protocol = "http"
	ProtocolHTTPS Protocol = "https"
	ProtocolTCP   Protocol = "tcp"
	ProtocolTLS   Protocol = "tls"
	ProtocolUnix  Protocol = "unix"
	ProtocolUDP   Protocol = "udp"

	// 子协议
	ProtocolHTTP4 Protocol = "http4"
	ProtocolHTTP6 Protocol = "http6"

	ProtocolHTTPS4 Protocol = "https4"
	ProtocolHTTPS6 Protocol = "https6"

	ProtocolTCP4 Protocol = "tcp4"
	ProtocolTCP6 Protocol = "tcp6"

	ProtocolTLS4 Protocol = "tls4"
	ProtocolTLS6 Protocol = "tls6"
)

func AllProtocols() []Protocol {
	return []Protocol{ProtocolHTTP, ProtocolHTTPS, ProtocolTCP, ProtocolTLS, ProtocolUnix, ProtocolUDP, ProtocolHTTP4, ProtocolHTTP6, ProtocolHTTPS4, ProtocolHTTPS6, ProtocolTCP4, ProtocolTCP6, ProtocolTLS4, ProtocolTLS6}
}
