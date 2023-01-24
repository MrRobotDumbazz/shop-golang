package delivery

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"shop/internal/service"
	"strings"
)

type key int

const (
	tokenCtxKey key = iota
)

func (h *Handler) ValidateJWT(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			handler.ServeHTTP(w, r)
		} else {
			token := authHeader[1]
			claims, err := h.services.Auth.ParseToken(token, service.AccessSecret)
			if err != nil {
				log.Println(err)
			}
			seller, err := h.services.Auth.ValidateToken(claims, false)
			if err != nil {
				log.Println(err)
			}
			h.services.ExpireToken(claims)
			ctx := context.WithValue(r.Context(), tokenCtxKey, seller.ID)
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

/*
func (h *Handler) ValidateJWT(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := AuthHeaderTokenExtractor(r)
		seller := 0
		claims, err := h.services.ParseToken(token, service.AccessSecret)
		if err != nil {
			log.Println(err)
			seller = 0
		}
		s, err := h.services.ValidateToken(claims, false)
		seller = s.ID
		ctx := context.WithValue(r.Context(), tokenCtxKey, seller)
		handler.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthHeaderTokenExtractor(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no JWT.
	}

	authHeaderParts := strings.Fields(authHeader)
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}
*/
