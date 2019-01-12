package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	server, err := StartServer()
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer server.Close()

	server.routes()
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
