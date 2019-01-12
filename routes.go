package main

import (
	"net/http"
)

func (s *Server) routes() {
	http.HandleFunc("/", s.home)
	http.HandleFunc("/antichistes/", s.docs)
	http.HandleFunc("/antichistes/random", s.randoms(5))
	http.HandleFunc("/antichistes/random/one", s.randoms(1))
	http.HandleFunc("/antichistes/vote", s.vote)
	http.HandleFunc("/antichistes/send", s.send)
}
