package delivery

import (
	"errors"
	"log"
	"net/http"
	"shop/internal/service"
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
		sellerid, ok := r.Context().Value(keySellerID).(int)
		if !ok {
			h.Errors(w, http.StatusInternalServerError, "Don't working context")
			return
		}
		authorization := true
		if sellerid == 0 {
			authorization = false
		}

		t, err := template.ParseFiles("templates/homepage.html")
		if err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		var products []models.Product
		if len(r.URL.Query()) == 0 {
			products, err = h.services.GetNewAllProducts()
			if err != nil {
				log.Print("err:delivery:homepage: GetNewAllProducts")
				h.Errors(w, http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			products, err = h.services.GetAllProductsBy(r.URL.Query())
			if err != nil {
				if errors.Is(err, service.ErrInvalidQueryRequest) {
					h.Errors(w, http.StatusNotFound, "Invalid query request")
					return
				}
				h.Errors(w, http.StatusInternalServerError, err.Error())
				return
			}
		}
		p := struct {
			Products      []models.Product
			Authorization bool
		}{
			Products:      products,
			Authorization: authorization,
		}
		if err = t.Execute(w, p); err != nil {
			log.Print(err)
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}
