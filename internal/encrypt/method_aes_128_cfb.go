package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

type AES128CFBMethod struct {
	iv    []byte
	block cipher.Block
}

func (this *AES128CFBMethod) Init(key, iv []byte) error {
	// 判断key是否为32长度
	l := len(key)
	if l > 16 {
		key = key[:16]
	} else if l < 16 {
		key = append(key, bytes.Repeat([]byte{' '}, 16-l)...)
	}

	// 判断iv长度
	l2 := len(iv)
	if l2 > aes.BlockSize {
		iv = iv[:aes.BlockSize]
	} else if l2 < aes.BlockSize {
		iv = append(iv, bytes.Repeat([]byte{' '}, aes.BlockSize-l2)...)
	}

	this.iv = iv

	// block
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	this.block = block

	return nil
}

func (this *AES128CFBMethod) Encrypt(src []byte) (dst []byte, err error) {
	if len(src) == 0 {
		return
	}

	defer func() {
		err = RecoverMethodPanic(recover())
	}()

	dst = make([]byte, len(src))
	encrypter := cipher.NewCFBEncrypter(this.block, this.iv)
	encrypter.XORKeyStream(dst, src)

	return
}

func (this *AES128CFBMethod) Decrypt(dst []byte) (src []byte, err error) {
	if len(dst) == 0 {
		return
	}

	defer func() {
		err = RecoverMethodPanic(recover())
	}()

	src = make([]byte, len(dst))
	encrypter := cipher.NewCFBDecrypter(this.block, this.iv)
	encrypter.XORKeyStream(src, dst)

	return
}
