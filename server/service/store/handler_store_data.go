package store

import (
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

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
	_, err = strconv.ParseUint(dealIDStr, 10, 64)
	if err != nil {
		log.Errorf("failed to parse deal ID: %s", err.Error())
		http.Error(w, "failed to parse deal ID", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("failed to read data: %v", err.Error())
		http.Error(w, "failed to read data", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, dealIDStr)

	_, err = os.Stat(path)
	if err != nil {
		//If the directory does not exist, make directory with dealID.
		if os.IsNotExist(err) {
			err = os.Mkdir(dealIDStr, os.ModePerm)
			if err != nil {
				log.Errorf("failed to create directory: %v", err.Error())
				http.Error(w, "failed to create directory", http.StatusInternalServerError)
				return
			}

			file, err := writeFile(path, dataHashStr, body)
			if err != nil {
				log.Errorf("failed to store data file: %v", err.Error())
				http.Error(w, "failed to store data file", http.StatusInternalServerError)
				return
			}
			defer file.Close()
		}
	}

	// If the directory exists, just write file with dataHash
	if os.IsExist(err) {
		file, err := writeFile(path, dataHashStr, body)
		if err != nil {
			log.Errorf("failed to store data file: %v", err.Error())
			http.Error(w, "failed to store data file", http.StatusInternalServerError)
			return
		}
		defer file.Close()
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/text")
	if _, err = w.Write([]byte(fmt.Sprintf("success to store data in %s/", path))); err != nil {
		log.Errorf("failed to write response: %s", err.Error())
		return
	}
}

func writeFile(path, dataHashStr string, body []byte) (*os.File, error) {
	path = filepath.Join(path, dataHashStr)

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %v", err.Error())
	}
	_, err = file.Write(body)
	if err != nil {
		return nil, fmt.Errorf("failed to write data: %v", err.Error())
	}

	// For overwrite file which has same filename(dataHash).
	_, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err.Error())
	}

	return file, nil
}
