package delivery

import (
	"log"
	"net/http"
	"shop/internal/service"
	"text/template"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	switch r.Method {
	case "GET":
		claims, ok := r.Context().Value(tokenCtxKey).(service.TokenClaims)
		if !ok {
			h.Errors(w, http.StatusInternalServerError, "Don't working context")
			return
		}
		seller, err := h.services.ValidateToken(claims, false)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		t, err := template.ParseFiles("templates/homepage.html")
		if err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err = t.Execute(w, seller); err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}
