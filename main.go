package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	db *sqlx.DB
)

// AntiJoke is the basic type for this API.
type AntiJoke struct {
	ID         int    `json:"id" db:"id"`
	FirstPart  string `json:"first_part" db:"first_part"`
	SecondPart string `json:"second_part" db:"second_part"`
}

func main() {
	host := os.Getenv("DB_HOST")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	connect := fmt.Sprintf("host=%s user=%s dbname=%s password=%s", host, user, name, pass)

	var err error
	db, err = sqlx.Connect("postgres", connect)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer db.Close()

	http.HandleFunc("/random", handleRandoms(5))
	http.HandleFunc("/random/one", handleRandoms(1))
	http.HandleFunc("/vote", handleVote)
	http.HandleFunc("/send", handleSend)

	err = http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func handleRandoms(limit int) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		var acs []AntiJoke
		err := db.Select(&acs, "SELECT id, first_part, second_part FROM antichistes WHERE public=true ORDER BY RANDOM() LIMIT $1", limit)
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("content-type", "application/json")
		json.NewEncoder(w).Encode(acs)
	}
}

func handleVote(w http.ResponseWriter, r *http.Request) {
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

	_, err = db.Exec("UPDATE antichistes SET VOTES=(SELECT votes FROM antichistes WHERE id=$1)+1 WHERE id=$1 AND public=true", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleSend(w http.ResponseWriter, r *http.Request) {
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

	_, err := db.Exec("INSERT INTO antichistes (first_part, second_part) VALUES ($1, $2)", firstPart, secondPart)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
