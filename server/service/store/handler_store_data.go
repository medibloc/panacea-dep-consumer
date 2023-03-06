package store

import (
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

	dealIDStr := mux.Vars(r)["dealId"]

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		log.Errorf("failed to read data: %v", err.Error())
		http.Error(w, "failed to read data", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	cwd, _ := os.Getwd()
	err = os.Mkdir(dealIDStr, os.ModePerm)
	if err != nil {
		log.Errorf(err.Error())
	}

	path := filepath.Join(cwd, dealIDStr, dataHashStr)
	newFilePath := filepath.FromSlash(path)
	_, err = os.OpenFile(newFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Errorf(err.Error())
	}

	file, err := os.Create(newFilePath)
	if err != nil {
		log.Errorf("failed to create file: %v", err.Error())
	}
	defer file.Close()

	_, err = file.Write(body)
	if err != nil {
		log.Errorf("failed to write data: %v", err.Error())
		return
	}
}
