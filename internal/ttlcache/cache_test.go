package ttlcache

import (
	"github.com/iwind/TeaGo/rands"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache := NewCache()
	cache.Write("a", 1, time.Now().Unix()+3600)
	cache.Write("b", 2, time.Now().Unix()+3601)
	cache.Write("a", 1, time.Now().Unix()+3602)
	cache.Write("d", 1, time.Now().Unix()+1)

	for _, piece := range cache.pieces {
		if len(piece.m) > 0 {
			for k, item := range piece.m {
				t.Log(k, "=>", item.Value, item.expiredAt)
			}
		}
	}
	t.Log(cache.Read("a"))
	time.Sleep(2 * time.Second)
	t.Log(cache.Read("d"))
}

func BenchmarkCache_Add(b *testing.B) {
	runtime.GOMAXPROCS(1)

	cache := NewCache()
	for i := 0; i < b.N; i++ {
		cache.Write(strconv.Itoa(i), i, time.Now().Unix()+int64(i%1024))
	}
}

func TestCache_IncreaseInt64(t *testing.T) {
	var cache = NewCache()

	{
		cache.IncreaseInt64("a", 1, time.Now().Unix()+3600)
		t.Log(cache.Read("a"))
	}
	{
		cache.IncreaseInt64("a", 1, time.Now().Unix()+3600+1)
		t.Log(cache.Read("a"))
	}
	{
		cache.Write("b", 1, time.Now().Unix()+3600+2)
		t.Log(cache.Read("b"))
	}
	{
		cache.IncreaseInt64("b", 1, time.Now().Unix()+3600+3)
		t.Log(cache.Read("b"))
	}
}

func TestCache_Read(t *testing.T) {
	runtime.GOMAXPROCS(1)

	var cache = NewCache(PiecesOption{Count: 32})

	for i := 0; i < 10_000_000; i++ {
		cache.Write("HELLO_WORLD_"+strconv.Itoa(i), i, time.Now().Unix()+int64(i%10240)+1)
	}

	total := 0
	for _, piece := range cache.pieces {
		//t.Log(len(piece.m), "keys")
		total += len(piece.m)
	}
	t.Log(total, "total keys")

	before := time.Now()
	for i := 0; i < 10_240; i++ {
		_ = cache.Read("HELLO_WORLD_" + strconv.Itoa(i))
	}
	t.Log(time.Since(before).Seconds()*1000, "ms")
}

func TestCache_GC(t *testing.T) {
	var cache = NewCache(&PiecesOption{Count: 5})
	cache.Write("a", 1, time.Now().Unix()+1)
	cache.Write("b", 2, time.Now().Unix()+2)
	cache.Write("c", 3, time.Now().Unix()+3)
	cache.Write("d", 4, time.Now().Unix()+4)
	cache.Write("e", 5, time.Now().Unix()+10)

	go func() {
		for i := 0; i < 1000; i++ {
			cache.Write("f", 1, time.Now().Unix()+1)
			time.Sleep(10 * time.Millisecond)
		}
	}()

	for i := 0; i < 20; i++ {
		cache.GC()
		t.Log("items:", cache.Count())
		time.Sleep(1 * time.Second)
	}

	t.Log("now:", time.Now().Unix())
	for _, p := range cache.pieces {
		for k, v := range p.m {
			t.Log(k, v.Value, v.expiredAt)
		}
	}
}

func TestCache_GC2(t *testing.T) {
	runtime.GOMAXPROCS(1)

	cache := NewCache()
	for i := 0; i < 1_000_000; i++ {
		cache.Write(strconv.Itoa(i), i, time.Now().Unix()+int64(rands.Int(0, 100)))
	}

	for i := 0; i < 100; i++ {
		t.Log(cache.Count(), "items")
		time.Sleep(1 * time.Second)
	}
}
