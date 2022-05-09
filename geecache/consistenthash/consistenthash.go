package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           // hash algorithm
	replicas int            // multiples of virtual nodes
	keys     []int          // sorted
	hashMap  map[int]string // mapping of virtual nodes to real nodes. key: position on the ring, value: name of real node
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) Add(peers ...string) {
	for _, peer := range peers {
		for i := 0; i < m.replicas; i++ {
			key := peer + "#" + strconv.Itoa(i)
			hash := int(m.hash([]byte(key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = peer
		}
	}
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

func (m *Map) Remove(peer string) {
	for i := 0; i < m.replicas; i++ {
		key := peer + "#" + strconv.Itoa(i)
		hash := int(m.hash([]byte(key)))
		delete(m.hashMap, hash)
		idx := sort.SearchInts(m.keys, hash)
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
	}
}

// todo: 此实现好像没有考虑到哈希冲突的问题，或者说哈希冲突了，就以最后一个存入的为准，这样 Remove 方法就有问题了
