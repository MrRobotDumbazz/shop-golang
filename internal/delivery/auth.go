package delivery

import (
	"errors"
	"log"
	"net/http"
	"shop/internal/service"
	"shop/models"
	"text/template"
	"time"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signup" {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("templates/signup.html")
		if err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err = t.Execute(w, nil); err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, err.Error())
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		email, ok := r.Form["email"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Please use latin letters")
			return
		}
		password, ok := r.Form["password"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Please use stronger password")
			return
		}
		seller := &models.Seller{
			Email:    email[0],
			Password: password[0],
		}
		if err := h.services.Auth.CreateSeller(seller); err != nil {
			if errors.Is(err, service.ErrInvalidEmail) || errors.Is(err, service.ErrInvalidPassword) {
				h.Errors(w, http.StatusBadRequest, err.Error())
				return
			}
			h.Error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/signin" {
		h.Errors(w, http.StatusNotFound, "")
	}
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("templates/signin.html")
		if err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err = t.Execute(w, nil); err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, err.Error())
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			h.Errors(w, http.StatusBadRequest, err.Error())
			return
		}
		email, ok := r.Form["email"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Please use latin letters")
			return
		}
		password, ok := r.Form["password"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Please use stronger password")
			return
		}
		token, err := h.services.GenerateJWT(email[0], password[0])
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "JWT",
			Value:   token,
			Path:    "/",
			Secure:  true,
			Expires: time.Now().Add(12 * time.Hour),
		})
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(tokenCtxKey).(service.TokenClaims)
	if !ok {
		return
	}
	h.services.DeleteToken(claims)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(tokenCtxKey).(service.TokenClaims)
	if !ok {
		return
	}
	user, err := h.services.ValidateToken(claims, true)
	if err != nil {
		h.Errors(w, http.StatusUnauthorized, "")
		return
	}
	_, err = h.services.GenerateRefreshJWT(user)
}
