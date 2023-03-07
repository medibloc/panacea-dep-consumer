package main

import (
	"flag"
	"os"

	"github.com/medibloc/panacea-dep-consumer/server"
)

func main() {
	httpPtr := flag.String("listen-addr", "", "http server listen address")
	grpcPtr := flag.String("grpc-addr", "", "grpc server listen address")
	chainIDPtr := flag.String("chain-id", "", "chain ID of Panacea blockchain")
	flag.Parse()
	if err := server.Run(*httpPtr, *grpcPtr, *chainIDPtr); err != nil {
		os.Exit(1)
	}
}
