package store_test

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/medibloc/panacea-dep-consumer/server/service/store"
	"github.com/stretchr/testify/require"
)

var (
	dataHash      string = "e2163e9619fe85fcab10f074b254e5901da9f37ec70f4d8f4539b927842cb58c"
	encryptedData string = "7D+aD626YZ9Q+0pLZO60G42nYr/rS4YrABbAckdAT6gxfNmLP1TgQ/hD6ZeqhAXVVGQw3pJzRZmYgj6ceU93zmShYroDTgv70+a+ZGdf6eRSIS0UipKGA9pREP5ZHKKIlUmoDvGTNzWWQR6HMh+eWiKiPUTJMAiQUchnPqcDhxU6moSF9TJJDeBkm4bNLreYG6blBWfckS5ZwQFQB63OTWb18YcVg+4v5Ho="
	dealId        string = "1"
)

func TestHandleStoreData(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	require.NoError(t, err)

	r, err := http.NewRequest("POST", "/v1/deals/{dealId}/data/{dataHash}", bytes.NewBuffer(data))
	require.NoError(t, err)

	r = mux.SetURLVars(r, map[string]string{"dealId": dealId, "dataHash": dataHash})

	w := httptest.NewRecorder()

	store.HandleStoreData(w, r)
	defer os.RemoveAll(dealId)

	require.Equal(t, http.StatusOK, w.Code)
}
