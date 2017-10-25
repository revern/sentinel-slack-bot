package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"strconv"
	"os"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if (err == nil) {
		log.Println("bad port")
	} else {
		stringPort := strconv.Itoa(port)
		router := mux.NewRouter().StrictSlash(true)
		router.HandleFunc("/hello/{name}", index).Methods("GET")
		http.ListenAndServe(":"+stringPort, router)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Println("Responsing to /hello request")
	log.Println(r.UserAgent())

	vars := mux.Vars(r)
	name := vars["name"]

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Hello:", name)
}
