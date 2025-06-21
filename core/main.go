package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"localp2p/api"
	"localp2p/config"
	"localp2p/discovery"
	"localp2p/transport"
)

func main() {
	var configPath = flag.String("config", "", "Path to config file")
	var rpcPort = flag.Int("rpc-port", 9090, "RPC server port")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Printf("Starting LocalP2P node: %s", cfg.NodeID)

	// Initialize components
	disc := discovery.NewDiscovery(cfg)
	trans := transport.NewTransport(cfg.NodeID, cfg.Port)
	rpc := api.NewRPCServer(disc, trans, *rpcPort)

	// Start discovery
	if err := disc.Start(); err != nil {
		log.Fatalf("Failed to start discovery: %v", err)
	}
	defer disc.Stop()

	// Start transport
	if err := trans.Start(); err != nil {
		log.Fatalf("Failed to start transport: %v", err)
	}
	defer trans.Stop()

	// Start RPC server
	go func() {
		if err := rpc.Start(); err != nil {
			log.Fatalf("Failed to start RPC server: %v", err)
		}
	}()

	// Start message handler
	go handleMessages(trans)

	// Wait for interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down...")
}

func handleMessages(trans *transport.Transport) {
	for msg := range trans.GetMessages() {
		log.Printf("Received message from %s: %s", msg.From, msg.Content)
	}
}
