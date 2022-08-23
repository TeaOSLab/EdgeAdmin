package csrf

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/TeaOSLab/EdgeAdmin/internal/configs"
	"github.com/iwind/TeaGo/types"
	"strconv"
	"time"
)

// Generate 生成Token
func Generate() string {
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	h := sha256.New()
	h.Write([]byte(configs.Secret))
	h.Write([]byte(timestamp))
	s := h.Sum(nil)
	token := base64.StdEncoding.EncodeToString([]byte(timestamp + fmt.Sprintf("%x", s)))
	sharedTokenManager.Put(token)
	return token
}

// Validate 校验Token
func Validate(token string) (b bool) {
	if len(token) == 0 {
		return
	}

	if !sharedTokenManager.Exists(token) {
		return
	}
	defer func() {
		sharedTokenManager.Delete(token)
	}()

	data, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return
	}

	hashString := string(data)
	if len(hashString) < 10+32 {
		return
	}

	timestampString := hashString[:10]
	hashString = hashString[10:]

	h := sha256.New()
	h.Write([]byte(configs.Secret))
	h.Write([]byte(timestampString))
	hashData := h.Sum(nil)
	if hashString != fmt.Sprintf("%x", hashData) {
		return
	}

	timestamp := types.Int64(timestampString)
	if timestamp < time.Now().Unix()-1800 { // 有效期半个小时
		return
	}

	return true
}
