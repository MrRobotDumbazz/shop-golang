package delivery

import (
	"log"
	"net/http"
	"shop/models"
	"text/template"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	switch r.Method {
	case "GET":
		claims, ok := r.Context().Value(tokenCtxKey).(int)
		if !ok {
			h.Errors(w, http.StatusInternalServerError, "Don't working context")
			return
		}
		seller := models.Seller{}
		if claims == 0 {
			seller = models.Seller{
				HasToken: false,
			}
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
