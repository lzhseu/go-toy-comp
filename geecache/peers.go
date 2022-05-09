package geecache

// PeerPicker an interface that must be implemented to
// locate the peer that owns a specific key
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

// PeerGetter an interface that must be implemented to Get the value by a peer
type PeerGetter interface {
	Get(group, key string) ([]byte, error)
}
