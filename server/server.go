package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-dep-consumer/server/service/store"
	log "github.com/sirupsen/logrus"
)

func Run(listenAddr string) error {
	router := mux.NewRouter()

	store.RegisterHandlers(router)

	server := &http.Server{
		Handler:      router,
		Addr:         listenAddr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Infof("HTTP server is started: %s", server.Addr)
	return server.ListenAndServe()
}
