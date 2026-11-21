package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/ArditZubaku/blockchain/node"
	"github.com/ArditZubaku/blockchain/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("Failed to listen: %+v", err)
	}

	opts := []grpc.ServerOption{}
	grpcServer := grpc.NewServer(opts...)
	node := node.NewNode()

	proto.RegisterNodeServer(grpcServer, node)

	go func() {
		for {
			time.Sleep(time.Second * 2)
			makeTransaction()
		}
	}()

	grpcServer.Serve(ln)
}

func makeTransaction() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	client, err := grpc.NewClient(":3000", opts...)
	if err != nil {
		log.Fatalf("Failed to dial: %+v", err)
	}

	nodeClient := proto.NewNodeClient(client)

	v := &proto.Version{Version: "0.1", Height: 1}

	// _, err = nodeClient.HandleTransaction(context.TODO(), &proto.Transaction{})
	_, err = nodeClient.Handshake(context.TODO(), v)
	if err != nil {
		// log.Fatalf("Failed to handle transaction: %+v", err)
		log.Fatalf("Failed to handshake: %+v", err)
	}
}
