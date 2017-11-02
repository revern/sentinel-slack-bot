package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"fmt"
	"github.com/gorilla/schema"
	"bytes"
	"github.com/robfig/cron"
	"time"
	"strconv"
	"net/url"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
	Cr     *cron.Cron
}

func (a *App) Initialize(dbUrl string) {
	//connectionString :=
	//	fmt.Sprintf("user=%s password=%s dbname=%s", user, password, dbname)

	//connection, _ := pq.ParseURL(dbUrl)
	//connection += " sslmode=require"
	var err error
	a.DB, err = sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal(err)
	}

	_, err = a.DB.Exec("CREATE TABLE IF NOT EXISTS " +
		`devices("name" varchar(50) PRIMARY KEY NOT NULL,` +
		`"location_id" varchar(50) NOT NULL,` +
		`"location" varchar(50) NOT NULL);`)

	a.Cr = cron.New()
	a.Cr.AddFunc("0 45 16 * * *", func() { a.remindToReturnDevices() })
	a.Cr.Start()

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
	a.Router.HandleFunc("/return", a.returnDevice).Methods("POST")
	a.Router.HandleFunc("/remove", a.deleteDevice).Methods("POST")
	a.Router.HandleFunc("/ping", a.handlePing).Methods("GET")
	a.Router.HandleFunc("/info", a.handleInfo).Methods("POST")
	a.Router.HandleFunc("/call", a.handleCall).Methods("POST")
	a.Router.HandleFunc("/time", a.getTime).Methods("GET")
}

func (a *App) handleCall(w http.ResponseWriter, r *http.Request) {
	msg, err := getSlackMessage(r)
	if err != nil {
		fmt.Println("Error decoding")
	}

	resp, err := http.Get("https://sentinel-api.herokuapp.com/api/v1/users")

	var users []user
	if err != nil {
		//respondWithError(w, http.StatusBadRequest, "invalid token")
	} else {
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&users); err != nil {
			//respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			//return
		}
		defer r.Body.Close()

	}

	displayName := ""
	if users != nil {
		for i := 0; i < len(users); i++ {
			if users[i].UID == msg.UserId {
				displayName = users[i].Name
				break
			}
		}
	}
	if displayName == "" {
		displayName = msg.UserName
	}

	apiUrl := "https://sentinel-api.herokuapp.com/api/v1/devices/call"
	http.PostForm(apiUrl, url.Values{"uid": {msg.Text}, "display_name": {displayName}})
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "trying to call device with id "+msg.Text)
}

func (a *App) getTime(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, strconv.Itoa(time.Now().Hour())+" : "+strconv.Itoa(time.Now().Minute()))
}

func (a *App) getDevices(w http.ResponseWriter, r *http.Request) {
	devices, err := getDevices(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	devicesInfo := ""
	for i := 0; i < len(devices); i++ {
		if devices[i].Location == "box" {
			devicesInfo += devices[i].Name + " location: box \n"
		} else {
			devicesInfo += devices[i].Name + " location: <@" + devices[i].LocationID + "> \n"
		}
	}

	fmt.Fprint(w, devicesInfo)
}

func (a *App) remindToReturnDevices() {
	devices, err := getDevices(a.DB)
	if err != nil {
		return
	}

	slack_msg := ""
	for i := 0; i < len(devices); i++ {
		if devices[i].Location != "box" {
			slack_msg += "<@" + devices[i].LocationID + "> please don't forget return < " + devices[i].Name + " > to the box\n"
		}
	}
	postSlackMessage(slack_msg)
}

func (a *App) takeDevice(w http.ResponseWriter, r *http.Request) {
	msg, err := getSlackMessage(r)
	if err != nil {
		fmt.Println("Error decoding")
	}

	d := device{Name: msg.Text, LocationID: msg.UserId, Location: msg.UserName}
	if err := d.updateDevice(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Device not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	postSlackMessage("<@" + d.LocationID + "> took " + d.Name)
	//respondWithJSON(w, http.StatusOK, d)
}

func (a *App) returnDevice(w http.ResponseWriter, r *http.Request) {
	msg, err := getSlackMessage(r)
	if err != nil {
		fmt.Println("Error decoding")
	}

	d := device{Name: msg.Text, LocationID: "box", Location: "box"}
	if err := d.updateDevice(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Device not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}
	postSlackMessage("<@" + msg.UserId + "> returned " + d.Name)

	//respondWithJSON(w, http.StatusOK, d)
}

func (a *App) addDevice(w http.ResponseWriter, r *http.Request) {
	msg, err := getSlackMessage(r)
	if err != nil {
		fmt.Println("Error decoding")
	}

	d := device{Name: msg.Text}

	if err := d.createDevice(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		postSlackMessage("New device < " + msg.Text + " > was added to collection")
	}
	//respondWithJSON(w, http.StatusCreated, d)
}

func (a *App) deleteDevice(w http.ResponseWriter, r *http.Request) {
	msg, err := getSlackMessage(r)
	if err != nil {
		fmt.Println("Error decoding")
	}

	p := device{Name: msg.Text}
	if err := p.deleteDevice(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	} else {
		postSlackMessage("Device < " + msg.Text + " > was returned to Box")
	}

	//respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func (a *App) handleInfo(w http.ResponseWriter, r *http.Request) {
	info := "/devices-location - Show all devices location" +
		"\n/add - Add new device to the collection" +
		"\n/delete - Remove device from the collection" +
		"\n/take - Take device" +
		"\n/return - Return device" +
		"\n/call - Call iOS device by item ID (numbers on sticker on device)" +
		"\n/info - Show all bot commands"
	fmt.Fprint(w, info)
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
	res.WriteHeader(http.StatusOK)
	fmt.Println(res, "good")
}

func getSlackMessage(r *http.Request) (*slack_message, error) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form")
	}

	msg := new(slack_message)
	decoder := schema.NewDecoder()

	err = decoder.Decode(msg, r.Form)

	if err != nil {
		fmt.Println("Error decoding")
		return nil, err
	} else {
		return msg, nil
	}
}

func postSlackMessage(message string) {
	webhook_url := "https://hooks.slack.com/services/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"
	webhook_msg := webhook_message{Text: message}
	jsonValue, _ := json.Marshal(webhook_msg);
	http.Post(webhook_url, "application/json", bytes.NewBuffer(jsonValue))
}