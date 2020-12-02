package ttlcache

import "github.com/cespare/xxhash"

func HashKey(key []byte) uint64 {
	return xxhash.Sum64(key)
}
