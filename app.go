package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"time"
	"fmt"
	"github.com/gorilla/schema"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(dbUrl string) {
	//connectionString :=
	//	fmt.Sprintf("user=%s password=%s dbname=%s", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	_, err = a.DB.Exec("CREATE TABLE IF NOT EXISTS " +
		`device("name" PRIMARY KEY,` +
		`"location" varchar(50) DEFAULT box)`)

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run(port string) {
	log.Fatal(http.ListenAndServe(":"+port, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/devices", a.getDevices).Methods("POST")
	a.Router.HandleFunc("/add", a.addDevice).Methods("POST")
	a.Router.HandleFunc("/take", a.takeDevice).Methods("POST")
	a.Router.HandleFunc("/take", a.returnDevice).Methods("POST")
	a.Router.HandleFunc("/remove", a.deleteDevice).Methods("POST")
	a.Router.HandleFunc("/ping", a.handlePing).Methods("GET")
}

func (a *App) getDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := getDevices(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, devices)
}

func (a *App) takeDevice(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		fmt.Println("Error parsing form")
	}

	msg := new(slack_message)
	decoder := schema.NewDecoder()

	err = decoder.Decode(msg, r.Form)

	if err != nil {
		fmt.Println("Error decoding")
	}

	d := device{Name: msg.Text}
	if err := d.updateDevice(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Device not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	fmt.Fprint(w, msg.UserName+" take "+msg.Text)

	//respondWithJSON(w, http.StatusOK, d)
}

//TODO rework
func (a *App) returnDevice(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form")
	}

	msg := new(slack_message)
	decoder := schema.NewDecoder()

	err = decoder.Decode(msg, r.Form)

	if err != nil {
		fmt.Println("Error decoding")
	}

	d := device{Name: msg.Text}
	if err := d.updateDevice(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Device not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	fmt.Fprint(w, msg.UserName+" take "+msg.Text)

	//respondWithJSON(w, http.StatusOK, d)
}

func (a *App) addDevice(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form")
	}

	msg := new(slack_message)
	decoder := schema.NewDecoder()

	err = decoder.Decode(msg, r.Form)

	if err != nil {
		fmt.Println("Error decoding")
	}

	d := device{Name: msg.Text}

	if err := d.createDevice(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Fprint(w, "New device < "+msg.Text+" > was added to collection")

	//respondWithJSON(w, http.StatusCreated, d)
}

func (a *App) deleteDevice(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//id, err := strconv.Atoi(vars["id"])
	//if err != nil {
	//	respondWithError(w, http.StatusBadRequest, "Invalid Product ID")
	//	return
	//}

	p := device{Name: "todo_remake"}
	if err := p.deleteDevice(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (a *App) handlePing(res http.ResponseWriter, req *http.Request) {
	time.Sleep(2 * time.Minute)
	pingSelf()
}

func pingSelf() {
	http.Get("https://whispering-ridge-24474.herokuapp.com/ping")
}
