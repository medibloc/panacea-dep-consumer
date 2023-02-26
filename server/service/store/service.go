package store

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/v1/data/{key}", HandleStoreData).Methods(http.MethodPost)
}
