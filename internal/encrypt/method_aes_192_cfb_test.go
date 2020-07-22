package encrypt

import (
	"runtime"
	"strings"
	"testing"
)

func TestAES192CFBMethod_Encrypt(t *testing.T) {
	method, err := NewMethodInstance("aes-192-cfb", "abc", "123")
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

func BenchmarkAES192CFBMethod_Encrypt(b *testing.B) {
	runtime.GOMAXPROCS(1)

	method, err := NewMethodInstance("aes-192-cfb", "abc", "123")
	if err != nil {
		b.Fatal(err)
	}

	src := []byte(strings.Repeat("Hello", 1024))
	for i := 0; i < b.N; i++ {
		dst, err := method.Encrypt(src)
		if err != nil {
			b.Fatal(err)
		}
		_ = dst
	}
}
