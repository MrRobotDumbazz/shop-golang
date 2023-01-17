package delivery

import (
	"net/http"
	"shop/internal/service"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(tokenCtxKey).(*service.TokenClaims)
	if !ok {
		return
	}
	claims, err := h.services.ParseToken(token, service.AccessSecret)
	if err != nil {
		return
	}
	h.services.DeleteToken(claims)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(tokenCtxKey).(*service.TokenClaims)
	if !ok {
		return
	}
	claims, err := h.services.ParseToken(token)
	if err != nil {
		return
	}
	err = h.services.ValidateToken(claims, true)
	if err != nil {
		h.Errors(w, http.StatusUnauthorized, "")
		return
	}
}
