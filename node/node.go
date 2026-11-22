package node

import (
	"context"
	"net"
	"slices"
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

func (n *Node) Start(bootstrapNodes []string) error {
	ln, err := net.Listen("tcp", n.listenAddr)
	if err != nil {
		return err
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	proto.RegisterNodeServer(grpcServer, n)

	zap.S().Infow("Node running on", "addr:", n.listenAddr)

	go func() {
		if err := n.bootstrapNetwork(bootstrapNodes); err != nil {
			zap.S().Fatalf("Failed to boostrap the network: %+v", err)
		}
	}()

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

func (n *Node) bootstrapNetwork(addrs []string) error {
	for _, addr := range addrs {
		if !n.canConnectWith(addr) {
			continue
		}

		zap.S().Debugw(
			"Dialing peer",
			"addr", n.listenAddr,
			"peer_addr", addr,
		)
		c, v, err := n.dialPeer(addr)
		if err != nil {
			return err
		}
		n.addPeer(c, v)
	}

	return nil
}

func (n *Node) dialPeer(addr string) (proto.NodeClient, *proto.Version, error) {
	c, err := makeNodeClient(addr)
	if err != nil {
		return nil, nil, err
	}

	v, err := c.Handshake(context.Background(), n.getVersion())
	if err != nil {
		zap.S().Errorf("Failed to handshake: %+v\n", err)
		return nil, nil, err
	}

	return c, v, nil
}

func (n *Node) addPeer(c proto.NodeClient, v *proto.Version) {
	n.peerLock.Lock()
	defer n.peerLock.Unlock()

	n.peers[c] = v

	go n.bootstrapNetwork(v.PeerList)

	zap.S().Debugw(
		"New peer successfully connected",
		"addr", n.listenAddr,
		"peer_addr", v.GetListenAddr(),
		"height", v.GetHeight(),
	)
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
		PeerList:   n.getPeerList(),
	}
}

func (n *Node) getPeerList() []string {
	n.peerLock.RLock()
	defer n.peerLock.RUnlock()

	peers := make([]string, 0, len(n.peers))
	for _, v := range n.peers {
		peers = append(peers, v.ListenAddr)
	}

	return peers
}

func (n *Node) canConnectWith(addr string) bool {
	if n.listenAddr == addr {
		return false
	}
	return !slices.Contains(n.getPeerList(), addr)
}

func makeNodeClient(listenAddr string) (proto.NodeClient, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	c, err := grpc.NewClient(listenAddr, opts...)
	if err != nil {
		return nil, err
	}

	return proto.NewNodeClient(c), nil
}
