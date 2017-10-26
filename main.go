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
	"time"
	_ "github.com/lib/pq"
	"database/sql"
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
	Token          string `schema:"token"`
	TeamId         string `schema:"team_id"`
	TeamDomain     string `schema:"team_domain"`
	EnterpriseId   string `schema:"enterprise_id"`
	EnterpriseName string `schema:"enterprise_name"`
	ChannelId      string `schema:"channel_id"`
	ChannelName    string `schema:"channel_name"`
	UserId         string `schema:"user_id"`
	UserName       string `schema:"user_name"`
	Command        string `schema:"command"`
	Text           string `schema:"text"`
	ResponseUrl    string `schema:"response_url"`
	TriggerId      string `schema:"trigger_id"`
}

func main() {
	//db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//_, err = db.Exec("CREATE TABLE IF NOT EXISTS " +
	//	`device("name" PRIMARY KEY,` +
	//	`"location" varchar(50))`)

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if (err != nil) {
		log.Println("bad port")
	} else {
		stringPort := strconv.Itoa(port)
		router := mux.NewRouter()
		router.HandleFunc("/movies", handleMovies).Methods("GET")
		router.HandleFunc("/movie/{imdbKey}", handleMovie).Methods("GET", "DELETE", "POST")
		router.HandleFunc("/ping", handlePing).Methods("GET")
		router.HandleFunc("/take", handleTake).Methods("POST")
		http.ListenAndServe(":" + stringPort, router)

		pingSelf()
	}
}

func handlePing(res http.ResponseWriter, req *http.Request) {
	time.Sleep(2 * time.Minute)
	pingSelf()
}

func pingSelf() {
	http.Get("https://whispering-ridge-24474.herokuapp.com/ping")
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

	fmt.Fprint(res, msg.Text+" fu")
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
