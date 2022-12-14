package delivery

import (
	"shop/internal/service"

	"github.com/gorilla/mux"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) Handlers() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", HomePage).Methods("GET")
	r.HandleFunc("/sign-up", SignUp).Methods("GET", "POST")
	r.HandleFunc("/signin", SignIn).Methods("GET", "POST")
	return r
}
