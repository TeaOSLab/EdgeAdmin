package encrypt

import "testing"

func TestMagicKeyEncode(t *testing.T) {
	dst := MagicKeyEncode([]byte("Hello,World"))
	t.Log("dst:", string(dst))

	src := MagicKeyDecode(dst)
	t.Log("src:", string(src))
}
