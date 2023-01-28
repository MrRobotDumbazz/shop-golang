package delivery

import (
	"context"
	"net/http"
	"shop/internal/service"
)

type key int

const (
	tokenCtxKey key = iota
)

func (h *Handler) ValidateJWT(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := service.TokenClaims{}
		cookie, err := r.Cookie("JWT")
		if err != nil {
			if err == http.ErrNoCookie {
				claims = service.TokenClaims{}
			}
			claims = service.TokenClaims{}
		} else {
			claims, _ = h.services.Auth.ParseToken(cookie.Value, service.AccessSecret)
		}
		ctx := context.WithValue(r.Context(), tokenCtxKey, claims)
		handler.ServeHTTP(w, r.WithContext(ctx))
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
