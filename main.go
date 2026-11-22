package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/ArditZubaku/blockchain/node"
	"github.com/ArditZubaku/blockchain/proto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// logger, err := zap.NewProduction()
	encCfg := zap.NewDevelopmentEncoderConfig()
	encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encCfg.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
	encCfg.EncodeCaller = zapcore.ShortCallerEncoder
	encCfg.TimeKey = ""

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encCfg),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)

	logger := zap.New(core, zap.AddCaller())

	// Replaces Zap's default no-op global logger
	undo := zap.ReplaceGlobals(logger)
	defer undo()
	defer func() {
		_ = zap.L().Sync()
	}()

	zap.S().Infow("The logger is ready")

	makeNode(":3000")
	time.Sleep(time.Second)
	makeNode(":4000", ":3000")
	time.Sleep(time.Second)
	makeNode(":5000", ":4000")

	// go func() {
	// for {
	// time.Sleep(time.Second * 2)
	// makeTransaction()
	// }
	// }()
	select {}
}

func makeNode(listenAddr string, bootstrapNodes ...string) *node.Node {
	n := node.NewNode(listenAddr)
	go n.Start(bootstrapNodes)
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
