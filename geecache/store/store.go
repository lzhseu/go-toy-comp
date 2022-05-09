package store

type Store interface {
	Get(key string) (Value, bool)
	Add(key string, value Value)
	Remove(key string)
	Len() int
}

type storage struct {
	maxBytes  int64
	nBytes    int64
	OnEvicted EvictedFunc // optional and executed when an entry is purged
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

type EvictedFunc func(key string, value Value)

func New(maxBytes int64, onEvicted EvictedFunc) *storage {
	return &storage{maxBytes: maxBytes, OnEvicted: onEvicted}
}
