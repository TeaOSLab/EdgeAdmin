package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

type AES256CFBMethod struct {
	block cipher.Block
	iv    []byte
}

func (this *AES256CFBMethod) Init(key, iv []byte) error {
	// 判断key是否为32长度
	l := len(key)
	if l > 32 {
		key = key[:32]
	} else if l < 32 {
		key = append(key, bytes.Repeat([]byte{' '}, 32-l)...)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	this.block = block

	// 判断iv长度
	l2 := len(iv)
	if l2 > aes.BlockSize {
		iv = iv[:aes.BlockSize]
	} else if l2 < aes.BlockSize {
		iv = append(iv, bytes.Repeat([]byte{' '}, aes.BlockSize-l2)...)
	}
	this.iv = iv

	return nil
}

func (this *AES256CFBMethod) Encrypt(src []byte) (dst []byte, err error) {
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

func (this *AES256CFBMethod) Decrypt(dst []byte) (src []byte, err error) {
	if len(dst) == 0 {
		return
	}

	defer func() {
		err = RecoverMethodPanic(recover())
	}()

	src = make([]byte, len(dst))
	decrypter := cipher.NewCFBDecrypter(this.block, this.iv)
	decrypter.XORKeyStream(src, dst)

	return
}
