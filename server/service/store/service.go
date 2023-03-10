package store

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Service struct {
	dataDir string
}

func NewService(dataDir string) *Service {
	return &Service{
		dataDir: dataDir,
	}
}

func (s *Service) RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/v0/deals/{dealId}/data/{dataHash}", s.HandleStoreData).Methods(http.MethodPost)
}
