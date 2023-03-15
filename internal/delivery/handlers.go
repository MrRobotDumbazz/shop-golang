package delivery

import (
	"shop/internal/service"

	"github.com/go-chi/chi"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Handlers() *chi.Mux {
	r := chi.NewRouter()
	r.Use(h.ValidateJWT)
	r.Get("/", h.HomePage)
	r.Get("/signup", h.SignUp)
	r.Post("/signup", h.SignUp)
	r.Get("/signin", h.SignIn)
	r.Post("/signin", h.SignIn)
	r.Get("/logout", h.Logout)
	r.Get("/createproduct", h.create_product)
	r.Post("/createproduct", h.create_product)
	r.Get("/product", h.product)
	return r
}
