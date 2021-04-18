package ttlcache

import (
	"github.com/iwind/TeaGo/types"
	"sync"
	"time"
)

type Piece struct {
	m        map[uint64]*Item
	maxItems int
	locker   sync.RWMutex
}

func NewPiece(maxItems int) *Piece {
	return &Piece{m: map[uint64]*Item{}, maxItems: maxItems}
}

func (this *Piece) Add(key uint64, item *Item) () {
	this.locker.Lock()
	if len(this.m) >= this.maxItems {
		this.locker.Unlock()
		return
	}
	this.m[key] = item
	this.locker.Unlock()
}

func (this *Piece) IncreaseInt64(key uint64, delta int64, expiredAt int64) (result int64) {
	this.locker.Lock()
	item, ok := this.m[key]
	if ok {
		result = types.Int64(item.Value) + delta
		item.Value = result
		item.expiredAt = expiredAt
	} else {
		if len(this.m) < this.maxItems {
			result = delta
			this.m[key] = &Item{
				Value:     delta,
				expiredAt: expiredAt,
			}
		}
	}
	this.locker.Unlock()
	return
}

func (this *Piece) Delete(key uint64) {
	this.locker.Lock()
	delete(this.m, key)
	this.locker.Unlock()
}

func (this *Piece) Read(key uint64) (item *Item) {
	this.locker.RLock()
	item = this.m[key]
	if item != nil && item.expiredAt < time.Now().Unix() {
		item = nil
	}
	this.locker.RUnlock()

	return
}

func (this *Piece) Count() (count int) {
	this.locker.RLock()
	count = len(this.m)
	this.locker.RUnlock()
	return
}

func (this *Piece) GC() {
	this.locker.Lock()
	timestamp := time.Now().Unix()
	for k, item := range this.m {
		if item.expiredAt <= timestamp {
			delete(this.m, k)
		}
	}
	this.locker.Unlock()
}

func (this *Piece) Destroy() {
	this.locker.Lock()
	this.m = nil
	this.locker.Unlock()
}
