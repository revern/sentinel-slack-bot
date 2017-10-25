package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"strconv"
	"os"
	"github.com/gorilla/schema"
)

// Movie Struct
type Movie struct {
	Title  string `json:"title"`
	Rating string `json:"rating"`
	Year   string `json:"year"`
}

var movies = map[string]*Movie{
	"tt0076759": &Movie{Title: "Star Wars: A New Hope", Rating: "8.7", Year: "1977"},
	"tt0082971": &Movie{Title: "Indiana Jones: Raiders of the Lost Ark", Rating: "8.6", Year: "1981"},
}

type Device struct {
	Name     string
	Location string
}

type SlackMessage struct {
	token           string
	team_id         string
	team_domain     string
	enterprise_id   string
	enterprise_name string
	channel_id      string
	channel_name    string
	user_id         string
	user_name       string
	command         string
	text            string
	response_url    string
	trigger_id      string
}

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if (err != nil) {
		log.Println("bad port")
	} else {
		stringPort := strconv.Itoa(port)
		router := mux.NewRouter()
		router.HandleFunc("/movies", handleMovies).Methods("GET")
		router.HandleFunc("/movie/{imdbKey}", handleMovie).Methods("GET", "DELETE", "POST")
		router.HandleFunc("/take", handleTake).Methods("POST")
		http.ListenAndServe(":"+stringPort, router)
	}
}

func handleTake(res http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()

	if err != nil {
		fmt.Println("Error parsing form")
	}

	msg := new(SlackMessage)
	decoder := schema.NewDecoder()

	err = decoder.Decode(msg, req.Form)

	if err != nil {
		fmt.Println("Error decoding")
	}

	fmt.Fprint(res, msg.text + " fu")
}

func handleMovie(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(req)
	imdbKey := vars["imdbKey"]

	switch req.Method {
	case "GET":
		movie, ok := movies[imdbKey]
		if !ok {
			res.WriteHeader(http.StatusNotFound)
			fmt.Fprint(res, string("Movie not found"))
		}
		outgoingJSON, error := json.Marshal(movie)
		if error != nil {
			log.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(res, string(outgoingJSON))
	case "DELETE":
		delete(movies, imdbKey)
		res.WriteHeader(http.StatusNoContent)
	case "POST":
		movie := new(Movie)
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&movie)
		if error != nil {
			log.Println(error.Error())
			http.Error(res, error.Error(), http.StatusInternalServerError)
			return
		}
		movies[imdbKey] = movie
		outgoingJSON, err := json.Marshal(movie)
		if err != nil {
			log.Println(error.Error())
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusCreated)
		fmt.Fprint(res, string(outgoingJSON))
	}
}

func handleMovies(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	outgoingJSON, error := json.Marshal(movies)
	if error != nil {
		log.Println(error.Error())
		http.Error(res, error.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(res, string(outgoingJSON))
}
