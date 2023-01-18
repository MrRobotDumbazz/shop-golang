package delivery

import (
	"encoding/json"
	"net/http"
	"shop/internal/service"
	"shop/models"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := &request{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		h.Error(w, r, http.StatusBadRequest, err)
		return
	}
	u := &models.Seller{
		Email:    req.Email,
		Password: req.Password,
	}
	if err := h.services.Auth.CreateSeller(u); err != nil {
		h.Error(w, r, http.StatusUnprocessableEntity, err)
		return
	}
	h.respond(w, r, http.StatusCreated, u)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value(tokenCtxKey).(string)
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
	token, ok := r.Context().Value(tokenCtxKey).(string)
	if !ok {
		return
	}
	claims, err := h.services.ParseToken(token, service.RefreshSecret)
	if err != nil {
		return
	}
	user, err := h.services.ValidateToken(claims, true)
	if err != nil {
		h.Errors(w, http.StatusUnauthorized, "")
		return
	}
	_, err = h.services.GenerateRefreshJWT(user)
}
