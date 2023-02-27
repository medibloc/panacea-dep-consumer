package middleware_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/btcsuite/btcd/btcec"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/medibloc/panacea-dep-consumer/panacea"
	"github.com/medibloc/panacea-dep-consumer/server/middleware"
	"github.com/stretchr/testify/require"
)

var (
	testOraclePrivKey, _ = btcec.NewPrivateKey(btcec.S256())
	testOraclePubKey     = testOraclePrivKey.PubKey()
)

func TestAuthSuccess(t *testing.T) {
	jwt := testGenerateJWT(t, testOraclePrivKey, 10*time.Second)
	testHTTPRequest(
		t,
		&mockGRPCClient{},
		fmt.Sprintf("Bearer %s", string(jwt)),
		http.StatusOK,
		"",
	)
}

func testGenerateJWT(t *testing.T, privKey *btcec.PrivateKey, expiration time.Duration) []byte {
	now := time.Now().Truncate(time.Second)
	token, err := jwt.NewBuilder().
		IssuedAt(now).
		NotBefore(now).
		Expiration(now.Add(expiration)).
		Build()
	require.NoError(t, err)

	signedJWT, err := jwt.Sign(token, jwt.WithKey(jwa.ES256K, privKey.ToECDSA()))
	require.NoError(t, err)

	return signedJWT
}

func testHTTPRequest(t *testing.T, grpcClient panacea.GRPCClient, authorizationHeader string, statusCode int, errMsg string) {
	req := httptest.NewRequest("GET", "http://test.com", nil)
	req.Header.Set("Authorization", authorizationHeader)

	w := httptest.NewRecorder()

	testHandler := middleware.NewJWTAuthMiddleware(grpcClient).Middleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}),
	)
	testHandler.ServeHTTP(w, req)

	resp := w.Result()
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	require.Equal(t, statusCode, resp.StatusCode)
	if errMsg != "" {
		require.Equal(t, errMsg+"\n", string(body))
	}
}

type mockGRPCClient struct {
}

func (c *mockGRPCClient) GetOraclePubKey(_ context.Context) (*btcec.PublicKey, error) {
	return testOraclePrivKey.PubKey(), nil
}

func (c *mockGRPCClient) GetCdc() *codec.ProtoCodec {
	return nil
}

func (c *mockGRPCClient) GetChainID() string {
	return ""
}

func (c *mockGRPCClient) BroadcastTx(_ []byte) (*tx.BroadcastTxResponse, error) {
	return nil, nil
}

func (c *mockGRPCClient) Close() error {
	return nil
}
