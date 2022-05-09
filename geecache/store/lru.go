package store

import (
	"container/list"
)

type LRUCache struct {
	*storage
	ll   *list.List
	dict map[string]*list.Element
}

func NewLRUCache(maxBytes int64, onEvicted EvictedFunc) *LRUCache {
	return &LRUCache{
		New(maxBytes, onEvicted),
		list.New(),
		make(map[string]*list.Element),
	}
}

func (L *LRUCache) Get(key string) (value Value, ok bool) {
	if ele, ok := L.dict[key]; ok {
		L.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (L *LRUCache) Add(key string, value Value) {
	if ele, ok := L.dict[key]; ok { // if existed, update
		L.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		L.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else { // not existed, add
		entry := &entry{key, value}
		ele := L.ll.PushFront(entry)
		L.dict[key] = ele
		L.nBytes += int64(len(key)) + int64(value.Len())
	}
	// remove overflow store
	for L.maxBytes != 0 && L.maxBytes < L.nBytes {
		L.RemoveOldest()
	}
}

func (L *LRUCache) Remove(key string) {
	if ele, ok := L.dict[key]; ok {
		L.ll.Remove(ele)
		delete(L.dict, key)
		kv := ele.Value.(*entry)
		L.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if L.OnEvicted != nil {
			L.OnEvicted(kv.key, kv.value)
		}
	}
}

func (L *LRUCache) Len() int {
	return L.ll.Len()
}

func (L *LRUCache) RemoveOldest() {
	ele := L.ll.Back()
	if ele != nil {
		kv := ele.Value.(*entry)
		L.ll.Remove(ele)
		delete(L.dict, kv.key)
		L.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if L.OnEvicted != nil {
			L.OnEvicted(kv.key, kv.value)
		}
	}
}
