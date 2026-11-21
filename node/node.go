package node

import (
	"context"
	"fmt"

	"github.com/ArditZubaku/blockchain/proto"
	"google.golang.org/grpc/peer"
)

type Node struct {
	// peers map[net.Addr]*grpc.ClientConn
	version string

	proto.UnimplementedNodeServer
}

func NewNode() *Node {
	return &Node{version: "0.1"}
}

func (n *Node) HandleTransaction(ctx context.Context, tx *proto.Transaction) (*proto.Ack, error) {
	peer, ok := peer.FromContext(ctx)
	fmt.Println("ok=", ok)
	fmt.Println("Received transaction from: ", peer)
	return nil, nil
}

func (n *Node) Handshake(ctx context.Context, v *proto.Version) (*proto.Version, error) {
	ourVersion := &proto.Version{
		Version: n.version,
		Height:  100,
	}

	p, ok := peer.FromContext(ctx)
	if ok {
		fmt.Printf("Received version from %+v: %s\n", p.Addr, v)
	}

	return ourVersion, nil
}
