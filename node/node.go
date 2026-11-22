package node

import (
	"context"
	"net"
	"sync"

	"github.com/ArditZubaku/blockchain/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/peer"
)

type Node struct {
	version    string
	listenAddr string

	// logger *zap.SugaredLogger

	peerLock sync.RWMutex
	peers    map[proto.NodeClient]*proto.Version

	proto.UnimplementedNodeServer
}

func NewNode(listenAddr string) *Node {
	return &Node{
		version:    "0.1",
		listenAddr: listenAddr,
		peers:      make(map[proto.NodeClient]*proto.Version),
	}
}

func (n *Node) Start() error {
	ln, err := net.Listen("tcp", n.listenAddr)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterNodeServer(grpcServer, n)

	zap.S().Infow("Node running on ->", "addr:", n.listenAddr)

	return grpcServer.Serve(ln)
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, ok := peer.FromContext(ctx)
	zap.S().Infow("Received transaction from", "ok", ok, "peer", peer)
	return nil, nil
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	c, err := makeNodeClient(v.ListenAddr)
	if err != nil {
		return nil, err
	}
	n.addPeer(c, v)

	return n.getVersion(), nil
}

func (n *Node) BootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		c, err := makeNodeClient(addr)
		if err != nil {
			return err
		}

		v, err := c.Handshake(context.Background(), n.getVersion())
		if err != nil {
			zap.S().Errorf("Failed to handshake: %+v\n", err)
			continue
		}

		n.addPeer(c, v)
	}

	return nil
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	zap.S().Debugw("New peer connected", "addr", n.listenAddr, "peer_addr", v.GetListenAddr(), "height", v.GetHeight())

	n.peers[c] = v
}

func (n *Node) rmPeer(c proto.NodeClient) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	delete(n.peers, c)
}

func (n *Node) getVersion() *proto.Version {
	return &proto.Version{
		Version:    "0.1",
		Height:     0,
		ListenAddr: n.listenAddr,
	}
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	c, err := grpc.NewClient(listenAddr, opts...)
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(c), nil
}
