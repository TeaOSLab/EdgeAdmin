package encrypt

import "testing"

func TestFindMethodInstance(t *testing.T) {
	t.Log(NewMethodInstance("a", "b", ""))
	t.Log(NewMethodInstance("aes-256-cfb", "123456", ""))
}
