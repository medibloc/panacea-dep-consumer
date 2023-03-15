package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/medibloc/panacea-dep-consumer/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	httpPtr := flag.String("listen-addr", "", "http server listen address")
	grpcPtr := flag.String("panacea-grpc-addr", "", "panacea grpc server listen address")
	dataDirPtr := flag.String("data-dir", "", "the path which data will be stored")
	flag.Parse()

	if *httpPtr == "" || *grpcPtr == "" || *dataDirPtr == "" {
		fmt.Fprintln(os.Stderr, "missing required flag")
		flag.Usage()
		os.Exit(1)
	}

	if err := server.Run(*httpPtr, *grpcPtr, *dataDirPtr); err != nil {
		log.Errorf("failed to start consumer service: %v", err)
		os.Exit(1)
	}

}
