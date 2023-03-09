package store

import (
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func HandleStoreData(w http.ResponseWriter, r *http.Request) {
	dataHashStr := mux.Vars(r)["dataHash"]
	_, err := hex.DecodeString(dataHashStr)
	if err != nil {
		log.Errorf("failed to decode dataHash: %s", err.Error())
		http.Error(w, "failed to decode dataHash", http.StatusBadRequest)
		return
	}

	dealIDStr := mux.Vars(r)["dealId"]
	_, err = strconv.ParseUint(dealIDStr, 10, 64)
	if err != nil {
		log.Errorf("failed to parse deal ID: %s", err.Error())
		http.Error(w, "failed to parse deal ID", http.StatusBadRequest)
		return
	}

	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, dealIDStr)
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Errorf("failed to create directory: %v", err.Error())
		http.Error(w, "failed to create directory", http.StatusInternalServerError)
		return
	}

	file, err := os.Create(filepath.Join(path, dataHashStr))
	if err != nil {
		log.Errorf("failed to create file: %v", err.Error())
		http.Error(w, "failed to create file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Errorf("failed to write file: %v", err.Error())
		http.Error(w, "failed to write file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	if _, err = w.Write([]byte("success to store data")); err != nil {
		log.Errorf("failed to write response: %v", err.Error())
		return
	}
}
