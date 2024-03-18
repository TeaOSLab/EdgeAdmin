package ttlcache

import "github.com/cespare/xxhash/v2"

func HashKey(key []byte) uint64 {
	return xxhash.Sum64(key)
}
