package delivery

import (
	"log"
	"net/http"
	"shop/internal/service"
	"shop/models"
	"strconv"
	"text/template"
)

func (h *Handler) create_product(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/createproduct" {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("templates/createpost.html")
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
		claims, ok := r.Context().Value(tokenCtxKey).(*service.TokenClaims)
		if !ok {
			h.Errors(w, http.StatusInternalServerError, "Don't working context")
			return
		}
		log.Printf("Claims: %v", claims)
		seller, err := h.services.ValidateToken(claims, false)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		if err := r.ParseForm(); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		name, ok := r.Form["name"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Wrong name product")
			return
		}
		company, ok := r.Form["company"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Wrong name product")
			return
		}
		description, ok := r.Form["description"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Wrong name product")
			return
		}
		price, ok := r.Form["price"]
		if !ok {
			h.Errors(w, http.StatusBadRequest, "Wrong name product")
			return
		}
		categories, ok := r.Form["categories"]
		p, err := strconv.ParseFloat(price[0], 64)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		product := &models.Product{
			SellerID:    seller.ID,
			Name:        name[0],
			Company:     company[0],
			Description: description[0],
			Price:       p,
			Category:    categories[0],
		}
		err = h.services.Product.CreateProduct(product)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
		t, err := template.ParseFiles("templates/createpost.html")
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
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
	}
}

func (h *Handler) product(w http.ResponseWriter, r *http.Request) {
	productid, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		h.Errors(w, http.StatusNotFound, "")
		return
	}
	claims, ok := r.Context().Value(tokenCtxKey).(*service.TokenClaims)
	if !ok {
		h.Errors(w, http.StatusInternalServerError, "Don't working context")
		return
	}
	log.Printf("Claims: %v", claims)
	seller, err := h.services.ValidateToken(claims, false)
	if err != nil {
		h.Errors(w, http.StatusInternalServerError, err.Error())
		return
	}
	authorization := true
	if seller.ID == 0 {
		authorization = false
		h.Errors(w, http.StatusForbidden, "")
		return
	}
	switch r.Method {
	case "GET":
		t, err := template.ParseFiles("templates/createpost.html")
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
		product, err := h.services.GetProductByProductID(productid)
		if err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
		}
		pageProduct := struct {
			Product       *models.Product
			Seller        int
			Authorization bool
		}{
			product,
			seller.ID,
			authorization,
		}
		if err := t.Execute(w, pageProduct); err != nil {
			h.Errors(w, http.StatusInternalServerError, err.Error())
			return
		}
	default:
		h.Errors(w, http.StatusMethodNotAllowed, "")
		return
	}
}
