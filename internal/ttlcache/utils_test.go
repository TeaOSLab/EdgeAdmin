package ttlcache

import (
	"runtime"
	"testing"
)

func BenchmarkHashKey(b *testing.B) {
	runtime.GOMAXPROCS(1)
	for i := 0; i < b.N; i++ {
		HashKey([]byte("HELLO,WORLDHELLO,WORLDHELLO,WORLDHELLO,WORLDHELLO,WORLDHELLO,WORLD"))
	}
}
