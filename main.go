package main

import (
	"context"
	"log"

	"github.com/ArditZubaku/blockchain/node"
	"github.com/ArditZubaku/blockchain/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// logger, err := zap.NewProduction()
	loggerCfg := zap.NewDevelopmentConfig()
	loggerCfg.EncoderConfig.TimeKey = ""
	logger, err := loggerCfg.Build()
	if err != nil {
		log.Fatal(err)
	}

	// Replaces Zap's default no-op global logger
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	defer func() {
		_ = zap.L().Sync()
	}()

	zap.S().Info("The logger is ready")

	makeNode(":3000")
	makeNode(":4000", ":3000")

	// go func() {
	// for {
	// time.Sleep(time.Second * 2)
	// makeTransaction()
	// }
	// }()
	select {}
}

func makeNode(listenAddr string, boostrapNodes ...string) *node.Node {
	n := node.NewNode(listenAddr)
	go n.Start()

	if err := n.BootstrapNetwork(boostrapNodes); err != nil {
		log.Fatalf("Failed to boostrap the network: %+v", err)
	}

	return n
}

func makeTransaction() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	client, err := grpc.NewClient(":3000", opts...)
	if err != nil {
		log.Fatalf("Failed to dial: %+v", err)
	}

	nodeClient := proto.NewNodeClient(client)

	v := &proto.Version{
		Version:    "0.1",
		Height:     1,
		ListenAddr: ":4000",
	}

	// _, err = nodeClient.HandleTransaction(context.TODO(), &proto.Transaction{})
	_, err = nodeClient.Handshake(context.TODO(), v)
	if err != nil {
		// log.Fatalf("Failed to handle transaction: %+v", err)
		log.Fatalf("Failed to handshake: %+v", err)
	}
}
