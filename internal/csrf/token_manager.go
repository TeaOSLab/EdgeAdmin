package csrf

import (
	"sync"
	"time"
)

var sharedTokenManager = NewTokenManager()

func init() {
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			sharedTokenManager.Clean()
		}
	}()
}

type TokenManager struct {
	tokenMap map[string]int64 // token => timestamp

	locker sync.Mutex
}

func NewTokenManager() *TokenManager {
	return &TokenManager{
		tokenMap: map[string]int64{},
	}
}

func (this *TokenManager) Put(token string) {
	this.locker.Lock()
	this.tokenMap[token] = time.Now().Unix()
	this.locker.Unlock()
}

func (this *TokenManager) Exists(token string) bool {
	this.locker.Lock()
	_, ok := this.tokenMap[token]
	this.locker.Unlock()
	return ok
}

func (this *TokenManager) Delete(token string) {
	this.locker.Lock()
	delete(this.tokenMap, token)
	this.locker.Unlock()
}

func (this *TokenManager) Clean() {
	this.locker.Lock()
	for token, timestamp := range this.tokenMap {
		if time.Now().Unix()-timestamp > 3600 { // 删除一个小时前的
			delete(this.tokenMap, token)
		}
	}
	this.locker.Unlock()
}
