package delivery

import (
	"context"
	"fmt"
	"net/http"
	"shop/internal/service"
	"strings"
)

type key int

const (
	tokenCtxKey key = iota
)

func (h *Handler) ValidateJWT(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Malformed token")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Malformed Token"))
		} else {
			token := authHeader[1]
			claims, err := h.services.Auth.ParseToken(token, service.AccessSecret)
			if err != nil {
				return
			}
			err = h.services.Auth.ValidateToken(claims, false)
			if err != nil {
				return
			}
			h.services.ExpireToken(claims)
			ctx := context.WithValue(r.Context(), tokenCtxKey, claims)
			handler.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
