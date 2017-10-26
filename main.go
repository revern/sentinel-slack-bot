package main

import (
	"log"
	"strconv"
	"os"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Println("bad port")
	}
	portStr := strconv.Itoa(port)
	a := App{}
	a.Initialize(os.Getenv("DATABASE_URL"))
	a.Run(portStr)
}

//func main() {
//	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	_, err = db.Exec("CREATE TABLE IF NOT EXISTS " +
//		`device("name" PRIMARY KEY,` +
//		`"location" varchar(50))`)
//
//	port, err := strconv.Atoi(os.Getenv("PORT"))
//	if (err != nil) {
//		log.Println("bad port")
//	} else {
//		stringPort := strconv.Itoa(port)
//		router := mux.NewRouter()
//		router.HandleFunc("/movies", handleMovies).Methods("GET")
//		router.HandleFunc("/movie/{imdbKey}", handleMovie).Methods("GET", "DELETE", "POST")
//		router.HandleFunc("/take", handleTake).Methods("POST")
//		http.ListenAndServe(":" + stringPort, router)
//
//		pingSelf()
//	}
//}
//
//
//func handleTake(res http.ResponseWriter, req *http.Request) {
//	err := req.ParseForm()
//
//	if err != nil {
//		fmt.Println("Error parsing form")
//	}
//
//	msg := new(SlackMessage)
//	decoder := schema.NewDecoder()
//
//	err = decoder.Decode(msg, req.Form)
//
//	if err != nil {
//		fmt.Println("Error decoding")
//	}
//
//	fmt.Fprint(res, msg.Text+" fu")
//}
//
//func handleMovie(res http.ResponseWriter, req *http.Request) {
//	res.Header().Set("Content-Type", "application/json")
//	vars := mux.Vars(req)
//	imdbKey := vars["imdbKey"]
//
//	switch req.Method {
//	case "GET":
//		movie, ok := movies[imdbKey]
//		if !ok {
//			res.WriteHeader(http.StatusNotFound)
//			fmt.Fprint(res, string("Movie not found"))
//		}
//		outgoingJSON, error := json.Marshal(movie)
//		if error != nil {
//			log.Println(error.Error())
//			http.Error(res, error.Error(), http.StatusInternalServerError)
//			return
//		}
//		fmt.Fprint(res, string(outgoingJSON))
//	case "DELETE":
//		delete(movies, imdbKey)
//		res.WriteHeader(http.StatusNoContent)
//	case "POST":
//		movie := new(Movie)
//		decoder := json.NewDecoder(req.Body)
//		error := decoder.Decode(&movie)
//		if error != nil {
//			log.Println(error.Error())
//			http.Error(res, error.Error(), http.StatusInternalServerError)
//			return
//		}
//		movies[imdbKey] = movie
//		outgoingJSON, err := json.Marshal(movie)
//		if err != nil {
//			log.Println(error.Error())
//			http.Error(res, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		res.WriteHeader(http.StatusCreated)
//		fmt.Fprint(res, string(outgoingJSON))
//	}
//}
//
//func handleMovies(res http.ResponseWriter, req *http.Request) {
//	res.Header().Set("Content-Type", "application/json")
//	outgoingJSON, error := json.Marshal(movies)
//	if error != nil {
//		log.Println(error.Error())
//		http.Error(res, error.Error(), http.StatusInternalServerError)
//		return
//	}
//	fmt.Fprint(res, string(outgoingJSON))
//}
