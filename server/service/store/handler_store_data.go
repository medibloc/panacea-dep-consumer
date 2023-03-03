package store

import (
	"encoding/hex"
	"io"
	"net/http"
	"os"

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
	defer r.Body.Close()
	if err != nil {
		log.Errorf("failed to read data: %v", err.Error())
		http.Error(w, "failed to read data", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(body); err != nil {
		log.Errorf("failed to write response: %s", err.Error())
		return
	}

	err = os.WriteFile(dataHashStr, body, 0644)
	if err != nil {
		log.Errorf("failed to write data: %s", err.Error())
		return
	}
}
