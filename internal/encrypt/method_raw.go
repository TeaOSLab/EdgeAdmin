package encrypt

type RawMethod struct {
}

func (this *RawMethod) Init(key, iv []byte) error {
	return nil
}

func (this *RawMethod) Encrypt(src []byte) (dst []byte, err error) {
	if len(src) == 0 {
		return
	}
	dst = make([]byte, len(src))
	copy(dst, src)
	return
}

func (this *RawMethod) Decrypt(dst []byte) (src []byte, err error) {
	if len(dst) == 0 {
		return
	}
	src = make([]byte, len(dst))
	copy(src, dst)
	return
}
