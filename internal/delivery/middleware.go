package delivery

import (
	"context"
	"log"
	"net/http"
	"shop/internal/service"
)

type key int

const (
	keySellerID key = iota
)

func (h *Handler) ValidateJWT(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := &service.TokenClaims{}
		sellerid := 0
		cookie, err := r.Cookie("JWT")
		if err != nil {
			if err == http.ErrNoCookie {
				claims = nil
				log.Println("Nil cookie")
			}
			log.Printf("Error in cookie: %v", err)
			claims = nil
		} else {
			claims, err = h.services.Auth.ParseToken(cookie.Value, service.AccessSecret)
			if err != nil {
				claims = nil
				log.Printf("Error: %v", err)
			}
			log.Printf("Claims: %v", claims)
			sellerid, err = h.services.ValidateToken(claims, false)
			log.Println("Seller id: %d", sellerid)
			if err != nil {
				log.Printf("Error: %v", err)
			}
		}
		ctx := context.WithValue(r.Context(), keySellerID, sellerid)
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
