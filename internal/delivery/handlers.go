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
	r.Get("/", h.HomePage)
	r.Get("/sign-up", h.SignUp)
	r.Post("/sign-up", h.SignUp)
	return r
}
