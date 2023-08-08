package ttlcache

import (
	"github.com/TeaOSLab/EdgeAdmin/internal/utils"
	"time"
)

var DefaultCache = NewCache()

// TTL缓存
// 最大的缓存时间为30 * 86400
// Piece数据结构：
//
//	    Piece1          |  Piece2 | Piece3 | ...
//	[ Item1, Item2, ... |   ...
//
// KeyMap列表数据结构
// { timestamp1 => [key1, key2, ...] }, ...
type Cache struct {
	isDestroyed bool
	pieces      []*Piece
	countPieces uint64
	maxItems    int

	gcPieceIndex int
	ticker       *utils.Ticker
}

func NewCache(opt ...OptionInterface) *Cache {
	countPieces := 128
	maxItems := 1_000_000
	for _, option := range opt {
		if option == nil {
			continue
		}
		switch o := option.(type) {
		case *PiecesOption:
			if o.Count > 0 {
				countPieces = o.Count
			}
		case *MaxItemsOption:
			if o.Count > 0 {
				maxItems = o.Count
			}
		}
	}

	cache := &Cache{
		countPieces: uint64(countPieces),
		maxItems:    maxItems,
	}

	for i := 0; i < countPieces; i++ {
		cache.pieces = append(cache.pieces, NewPiece(maxItems/countPieces))
	}

	// start timer
	go func() {
		cache.ticker = utils.NewTicker(5 * time.Second)
		for cache.ticker.Next() {
			cache.GC()
		}
	}()

	return cache
}

func (this *Cache) Write(key string, value interface{}, expiredAt int64) {
	if this.isDestroyed {
		return
	}

	currentTimestamp := time.Now().Unix()
	if expiredAt <= currentTimestamp {
		return
	}

	maxExpiredAt := currentTimestamp + 30*86400
	if expiredAt > maxExpiredAt {
		expiredAt = maxExpiredAt
	}
	uint64Key := HashKey([]byte(key))
	pieceIndex := uint64Key % this.countPieces
	this.pieces[pieceIndex].Add(uint64Key, &Item{
		Value:     value,
		expiredAt: expiredAt,
	})
}

func (this *Cache) IncreaseInt64(key string, delta int64, expiredAt int64) int64 {
	if this.isDestroyed {
		return 0
	}

	currentTimestamp := time.Now().Unix()
	if expiredAt <= currentTimestamp {
		return 0
	}

	maxExpiredAt := currentTimestamp + 30*86400
	if expiredAt > maxExpiredAt {
		expiredAt = maxExpiredAt
	}
	uint64Key := HashKey([]byte(key))
	pieceIndex := uint64Key % this.countPieces
	return this.pieces[pieceIndex].IncreaseInt64(uint64Key, delta, expiredAt)
}

func (this *Cache) Read(key string) (item *Item) {
	uint64Key := HashKey([]byte(key))
	return this.pieces[uint64Key%this.countPieces].Read(uint64Key)
}

func (this *Cache) Delete(key string) {
	uint64Key := HashKey([]byte(key))
	this.pieces[uint64Key%this.countPieces].Delete(uint64Key)
}

func (this *Cache) Count() (count int) {
	for _, piece := range this.pieces {
		count += piece.Count()
	}
	return
}

func (this *Cache) GC() {
	this.pieces[this.gcPieceIndex].GC()
	newIndex := this.gcPieceIndex + 1
	if newIndex >= int(this.countPieces) {
		newIndex = 0
	}
	this.gcPieceIndex = newIndex
}

func (this *Cache) Destroy() {
	this.isDestroyed = true

	if this.ticker != nil {
		this.ticker.Stop()
		this.ticker = nil
	}
	for _, piece := range this.pieces {
		piece.Destroy()
	}
}
