package ttlcache

type Item struct {
	Value     interface{}
	expiredAt int64
}
