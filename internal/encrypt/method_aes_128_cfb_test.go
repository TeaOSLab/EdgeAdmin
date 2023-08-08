package encrypt

import (
	"runtime"
	"strings"
	"testing"
)

func TestAES128CFBMethod_Encrypt(t *testing.T) {
	method, err := NewMethodInstance("aes-128-cfb", "abc", "123")
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

func TestAES128CFBMethod_Encrypt2(t *testing.T) {
	method, err := NewMethodInstance("aes-128-cfb", "abc", "123")
	if err != nil {
		t.Fatal(err)
	}

	sources := [][]byte{}

	{
		a := []byte{1}
		_, err = method.Encrypt(a)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 10; i++ {
		src := []byte(strings.Repeat("Hello", 1))
		dst, err := method.Encrypt(src)
		if err != nil {
			t.Fatal(err)
		}

		sources = append(sources, dst)
	}

	{

		a := []byte{1}
		_, err = method.Decrypt(a)
		if err != nil {
			t.Fatal(err)
		}
	}

	for _, dst := range sources {
		dst2 := append([]byte{}, dst...)
		src2, err := method.Decrypt(dst2)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(string(src2))
	}
}

func BenchmarkAES128CFBMethod_Encrypt(b *testing.B) {
	runtime.GOMAXPROCS(1)

	method, err := NewMethodInstance("aes-128-cfb", "abc", "123")
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
