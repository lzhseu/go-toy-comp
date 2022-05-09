package geecache

import (
	"fmt"
	pb "geecache/geecachepb"
	"geecache/singleflight"
	"log"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (fn GetterFunc) Get(key string) ([]byte, error) {
	return fn(key)
}

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name       string
	getter     Getter
	mainCache  cache
	peerPicker PeerPicker
	loader     *singleflight.Group
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("Gee Cache: nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{strategy: LRU, cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

func (g *Group) RegisterPeerPicker(peerPicker PeerPicker) {
	if g.peerPicker != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peerPicker = peerPicker
}

// SetStrategy will not take effect after call cache.add
func (g *Group) SetStrategy(strategy string) {
	g.mainCache.strategy = strategy
}

func (g *Group) Get(key string) (bv ByteView, err error) {
	if key == "" {
		err = fmt.Errorf("key is required")
		return
	}

	if bv, ok := g.mainCache.get(key); ok {
		log.Println("[GeeCache] hit")
		return bv, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	bv, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peerPicker != nil {
			if peer, ok := g.peerPicker.PickPeer(key); ok {
				value, err := g.getFromPeer(peer, key)
				if err == nil {
					return value, nil
				}
				log.Println("[GeeCache] Fail to get from peer", err)
			}
		}
		return g.getLocally(key)
	})

	if err == nil {
		return bv.(ByteView), err
	}
	return ByteView{}, err
}

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	resp := &pb.Response{}
	err := peer.Get(req, resp)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{resp.Value}, nil
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	bv := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, bv)
	return bv, nil
}

func (g *Group) populateCache(key string, bv ByteView) {
	g.mainCache.add(key, bv)
}
