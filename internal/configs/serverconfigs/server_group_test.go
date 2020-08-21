package serverconfigs

import (
	"github.com/iwind/TeaGo/assert"
	"testing"
)

func TestServerGroup_Protocol(t *testing.T) {
	a := assert.NewAssertion(t)

	{
		group := NewServerGroup("tcp://127.0.0.1:1234")
		a.IsTrue(group.Protocol() == ProtocolTCP)
		a.IsTrue(group.Addr() == "127.0.0.1:1234")
	}

	{
		group := NewServerGroup("http4://127.0.0.1:1234")
		a.IsTrue(group.Protocol() == ProtocolHTTP4)
		a.IsTrue(group.Addr() == "127.0.0.1:1234")
	}

	{
		group := NewServerGroup("127.0.0.1:1234")
		a.IsTrue(group.Protocol() == ProtocolHTTP)
		a.IsTrue(group.Addr() == "127.0.0.1:1234")
	}

	{
		group := NewServerGroup("unix:/tmp/my.sock")
		a.IsTrue(group.Protocol() == ProtocolUnix)
		a.IsTrue(group.Addr() == "/tmp/my.sock")
	}
}
