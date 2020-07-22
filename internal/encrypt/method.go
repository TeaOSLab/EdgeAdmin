package encrypt

type MethodInterface interface {
	// 初始化
	Init(key []byte, iv []byte) error

	// 加密
	Encrypt(src []byte) (dst []byte, err error)

	// 解密
	Decrypt(dst []byte) (src []byte, err error)
}
