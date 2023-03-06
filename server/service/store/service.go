package store

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/v1/deals/{dealId}/data/{key}", HandleStoreData).Methods(http.MethodPost)
}
