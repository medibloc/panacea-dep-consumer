package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/lestrrat-go/jwx/v2/jwt"
)

type JwtAuthMiddleware struct {
	// TODO: manage a nonce per account
}

func NewJWTAuthMiddleware() *JwtAuthMiddleware {
	return &JwtAuthMiddleware{}
}

func (j *JwtAuthMiddleware) Middleware(next http.Handler) http.Handler {
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
		parsedJWT, err := jwt.ParseInsecure(jwtBz)
		if err != nil {
			http.Error(w, "invalid jwt", http.StatusUnauthorized)
			return
		}

		// pass the authenticated account address to next handlers
		newReq := r.WithContext(
			context.WithValue(r.Context(), ContextKeyAuthenticatedAccountAddress{}, parsedJWT.Issuer()),
		)

		next.ServeHTTP(w, newReq)
	})
}

func parseBearerToken(authHeader string) (string, error) {
	elems := strings.Split(authHeader, " ")
	if len(elems) != 2 || elems[0] != "Bearer" {
		return "", fmt.Errorf("invalid bearer token")
	}

	return elems[1], nil
}
