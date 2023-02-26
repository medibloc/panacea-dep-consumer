package store

import (
	"encoding/hex"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func HandleStoreData(w http.ResponseWriter, r *http.Request) {
	dataHashStr := mux.Vars(r)["key"]
	_, err := hex.DecodeString(dataHashStr)
	if err != nil {
		log.Errorf("failed to decode dataHash: %s", err.Error())
		http.Error(w, "failed to decode dataHash", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	if err != nil {
		log.Errorf("failed to read data: %s", err.Error())
		http.Error(w, "failed to read data", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(body)
	if err != nil {
		log.Errorf("failed to write response: %s", err.Error())
		return
	}
}
