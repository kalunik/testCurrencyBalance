package api

import (
	"github.com/go-chi/chi/v5"
)

type Router struct {
	Mux *chi.Mux
}

func NewRouter() *Router {
	return &Router{Mux: chi.NewRouter()}
}

func (r *Router) PathMetaRoutes(h Handlers) {
	r.Mux.Route("/api/v1", func(r chi.Router) {

		r.Post("/{id}/deposit", h.depositFunds)
		r.Post("/{id}/withdrawal", h.withdrawFunds)

		r.Get("/{id}/balance", h.getBalance)
	})
}
