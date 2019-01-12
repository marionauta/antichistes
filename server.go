package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Server struct {
	db *sqlx.DB
}

func StartServer() (*Server, error) {
	info, err := url.Parse(os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, err
	}

	host := info.Hostname()
	name := strings.TrimPrefix(info.Path, "/")
	user := info.User.Username()
	pass, _ := info.User.Password()
	connect := fmt.Sprintf("host=%s user=%s dbname=%s password=%s", host, user, name, pass)

	db, err := sqlx.Connect("postgres", connect)
	if err != nil {
		return nil, err
	}

	return &Server{db}, nil
}

func (s *Server) Close() {
	s.db.Close()
}

func (s *Server) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.Header().Set("location", "/")
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, "pages/index.html")
}

func (s *Server) docs(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "pages/docs.html")
}

func (s *Server) randoms(limit int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var acs []AntiJoke
		err := s.db.Select(&acs, "SELECT id, first_part, second_part FROM antichistes WHERE public=true ORDER BY RANDOM() LIMIT $1", limit)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res := APIResponse{
			Error: 0,
			Items: acs,
		}

		w.Header().Set("access-control-allow-origin", "*")
		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(res)
	}
}

func (s *Server) vote(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	id, err := strconv.ParseInt(r.PostFormValue("id"), 10, 0)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err = s.db.Exec("UPDATE antichistes SET VOTES=(SELECT votes FROM antichistes WHERE id=$1)+1 WHERE id=$1 AND public=true", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("access-control-allow-origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{Error: 0})
}

func (s *Server) send(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	firstPart := r.PostFormValue("first_part")
	secondPart := r.PostFormValue("second_part")
	if len(firstPart)+len(secondPart) < 10 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "too short")
		return
	}

	_, err := s.db.Exec("INSERT INTO antichistes (first_part, second_part) VALUES ($1, $2)", firstPart, secondPart)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("access-control-allow-origin", "*")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{Error: 0})
}
