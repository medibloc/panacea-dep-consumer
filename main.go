package main

import (
	"flag"
	"os"

	"github.com/medibloc/panacea-dep-consumer/server"
)

func main() {
	httpPtr := flag.String("listenAddr", "", "http server listen address")
	grpcPtr := flag.String("grpcAddr", "", "grpc server listen address")
	chainIDPtr := flag.String("chainID", "", "chain ID")
	flag.Parse()
	if err := server.Run(*httpPtr, *grpcPtr, *chainIDPtr); err != nil {
		os.Exit(1)
	}
}
