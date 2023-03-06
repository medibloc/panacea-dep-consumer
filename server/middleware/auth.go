package middleware

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/medibloc/panacea-dep-consumer/panacea"
	log "github.com/sirupsen/logrus"
)

type JwtAuthMiddleware struct {
	panaceaGRPCClient panacea.GRPCClient
}

func NewJWTAuthMiddleware(grpcClient panacea.GRPCClient) *JwtAuthMiddleware {
	return &JwtAuthMiddleware{
		panaceaGRPCClient: grpcClient,
	}
}

func (mw *JwtAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		jwtStr, err := parseBearerToken(authHeader)
		if err != nil {
			http.Error(w, "invalid bearer token", http.StatusUnauthorized)
			return
		}
		jwtBz := []byte(jwtStr)

		// parse JWT without signature verification to get payloads for retrieving an auth pubkey
		_, err = jwt.ParseInsecure(jwtBz)
		if err != nil {
			http.Error(w, "invalid jwt", http.StatusUnauthorized)
			return
		}

		oraclePubKey, err := mw.queryOracleParams(r.Context())
		if err != nil {
			log.Error(err)
			http.Error(w, "cannot query oracle pubkey", http.StatusUnauthorized)
			return
		}

		_, err = jwt.Parse(jwtBz, jwt.WithKey(jwa.ES256K, oraclePubKey))
		if err != nil {
			http.Error(w, "jwt signature verification failed", http.StatusUnauthorized)
			return
		}

		newReq := r.WithContext(
			context.WithValue(r.Context(), ContextOraclePubKey{}, oraclePubKey),
		)

		next.ServeHTTP(w, newReq)
	})
}

func (mw *JwtAuthMiddleware) queryOracleParams(ctx context.Context) (*ecdsa.PublicKey, error) {
	oraclePubKey, err := mw.panaceaGRPCClient.GetOraclePubKey(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query oracle pubkey: %w", err)
	}

	return oraclePubKey.ToECDSA(), nil
}

func parseBearerToken(authHeader string) (string, error) {
	elems := strings.Split(authHeader, " ")
	if len(elems) != 2 || elems[0] != "Bearer" {
		return "", fmt.Errorf("invalid bearer token")
	}

	return elems[1], nil
}
