package main

import (
	"flag"
	"os"

	"github.com/medibloc/panacea-dep-consumer/server"
)

func main() {
	flag.String("listenAddr")
	if err := server.Run("127.0.0.1:8080"); err != nil {
		os.Exit(1)
	}
}
