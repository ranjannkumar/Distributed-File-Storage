package p2p





// HanshakeFUnc...?
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error { return nil }
