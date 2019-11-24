// app.go

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang/gddo/httputil/header"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(host, port, user, password, dbname, sslmode string) {
	connectionString :=
		fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	a.resetTables(a.DB)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	fmt.Println("DB and router initialized")
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}

func (a *App) initializeRoutes() {
	searchAPI := a.Router.PathPrefix("/search").Subrouter()
	updateAPI := a.Router.PathPrefix("/update").Subrouter()
	updateAPI.HandleFunc("/rack", a.createRack).Methods("POST")
	updateAPI.HandleFunc("/pdu", a.createPDU).Methods("POST")
	searchAPI.HandleFunc("/rack/{id:[0-9]+}", a.getRack).Methods("GET")
	searchAPI.HandleFunc("/pdu/{id:[0-9]+}", a.getPDU).Methods("GET")
	updateAPI.HandleFunc("/rack/{id:[0-9]+}", a.updateRack).Methods("PUT")
	updateAPI.HandleFunc("/pdu/{id:[0-9]+}", a.updatePDU).Methods("PUT")
}

func (a *App) createRack(w http.ResponseWriter, r *http.Request) {
	if !ValidNonJSONRequest(w, r) {
		return
	}
	var p rack
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Use format '/update/rack'")
		return
	}
	//	query := r.URL.Query()
	//	p.Height, _ = strconv.Atoi(query.Get("height"))
	//	p.Width, _ = strconv.Atoi(query.Get("width"))
	//	p.Location = query.Get("location")
	defer r.Body.Close()
	if err := p.createRack(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) getRack(w http.ResponseWriter, r *http.Request) {
	if !ValidNonJSONRequest(w, r) {
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid rack ID")
		return
	}

	p := rack{ID: id}
	if err := p.getRack(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Object not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) updateRack(w http.ResponseWriter, r *http.Request) {
	if !ValidNonJSONRequest(w, r) {
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid rack ID")
		return
	}
	var p rack
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request, use format /update/rack/{id}")
		return
	}
	defer r.Body.Close()
	p.ID = id
	if err := p.updateRack(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) createPDU(w http.ResponseWriter, r *http.Request) {
	if !ValidNonJSONRequest(w, r) {
		return
	}
	var p pdu
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Use format '/update/pdu'")
		return
	}
	query := r.URL.Query()
	p.OutletCount, _ = strconv.Atoi(query.Get("outletcount"))
	p.OutletUsed, _ = strconv.Atoi(query.Get("outletused"))
	p.PowerCapacity, _ = strconv.Atoi(query.Get("powercapacity"))
	defer r.Body.Close()
	if err := p.createPDU(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) getPDU(w http.ResponseWriter, r *http.Request) {
	if !ValidNonJSONRequest(w, r) {
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid pdu ID")
		return
	}

	p := pdu{ID: id}
	if err := p.getPDU(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Object not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) updatePDU(w http.ResponseWriter, r *http.Request) {
	if !ValidNonJSONRequest(w, r) {
		return
	}
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid pdu ID")
		return
	}
	var p pdu
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request, use format /update/pdu/{id}")
		return
	}
	defer r.Body.Close()
	query := r.URL.Query()
	p.ID = id
	p.OutletCount, _ = strconv.Atoi(query.Get("outletcount"))
	p.OutletUsed, _ = strconv.Atoi(query.Get("outletused"))
	p.PowerCapacity, _ = strconv.Atoi(query.Get("powercapacity"))
	if err := p.updatePDU(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
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

func ValidNonJSONRequest(w http.ResponseWriter, req *http.Request) bool {

	if req.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(req.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			respondWithError(w, http.StatusUnsupportedMediaType, msg)
			return false
		}
		return true
	}
	respondWithError(w, http.StatusBadRequest, "No content type")
	return false
}

func (a *App) resetTables(db *sql.DB) {
	a.DB.Exec(`DROP TABLE IF EXISTS racks;
	CREATE TABLE racks
	(
		id SERIAL,
		height INT NOT NULL,
		width INT NOT NULL,
		location VARCHAR(500) NOT NULL,
		CONSTRAINT racks_pkey PRIMARY KEY (id)
	)`)
	a.DB.Exec(`DROP TABLE IF EXISTS pdus
	CREATE TABLE pdus
	(
		id SERIAL,
		outletcount INT NOT NULL,
		outletused INT NOT NULL,
		powercapacity VARCHAR(500) NOT NULL,
		CONSTRAINT pdus_pkey PRIMARY KEY (id)
	)`)

}
