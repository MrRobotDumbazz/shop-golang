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
			return
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
			h.Errors(w, http.StatusUnprocessableEntity, err.Error())
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
		log.Printf("Token UUID: %s", token)
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
	sellerid, ok := r.Context().Value(keySellerID).(int)
	if !ok {
		h.Errors(w, http.StatusInternalServerError, "Don't working context")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "",
		Value:  "",
		Path:   "/",
		Secure: true,
	})
	h.services.DeleteToken(sellerid)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	sellerid, ok := r.Context().Value(keySellerID).(int)
	if !ok {
		h.Errors(w, http.StatusInternalServerError, "Don't working context")
		return
	}
	refreshtoken, err := h.services.GenerateRefreshJWT(sellerid)
	if err != nil {
		h.Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:   "JWT",
		Value:  refreshtoken,
		Path:   "/",
		Secure: true,
	})
}
