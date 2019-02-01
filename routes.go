package main

import (
	"github.com/go-chi/chi"
)

func (s *Server) routes() {
	s.r.Get("/", s.home)
	s.r.Route("/antichistes", func(r chi.Router) {
		r.Get("/", s.docs)
		r.Route("/random", func(r chi.Router) {
			r.Get("/", s.randoms(5))
			r.Get("/one", s.randoms(1))
		})
		r.Post("/vote", s.vote)
		r.Post("/send", s.send)
	})
}
