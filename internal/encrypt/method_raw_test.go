package encrypt

import "testing"

func TestRawMethod_Encrypt(t *testing.T) {
	method, err := NewMethodInstance("raw", "abc", "123")
	if err != nil {
		t.Fatal(err)
	}
	src := []byte("Hello, World")
	dst, err := method.Encrypt(src)
	if err != nil {
		t.Fatal(err)
	}
	dst = dst[:len(src)]
	t.Log("dst:", string(dst))

	src, err = method.Decrypt(dst)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("src:", string(src))
}
