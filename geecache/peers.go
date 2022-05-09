package geecache

import pb "geecache/geecachepb"

// PeerPicker an interface that must be implemented to
// locate the peer that owns a specific key
type PeerPicker interface {
	PickPeer(key string) (PeerGetter, bool)
}

// PeerGetter an interface that must be implemented to Get the value by a peer
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
