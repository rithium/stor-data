package main

import (
	"github.com/gorilla/mux"
	"github.com/rithium/stor-data/config"
	"log"
	"github.com/urfave/negroni"
	"net/http"
	"time"
	"fmt"
	"flag"
	"os"
	"github.com/rithium/version"
	"github.com/rithium/stor-data/model"
	"encoding/json"
	"strconv"
	"github.com/rithium/logger"
)

type Env struct {
	db model.Datastore
}

var rotatingLog *logger.RotatingFileWriter

func init() {
	versionFlag := flag.Bool("v", false, "prints version")
	configFlag := flag.Bool("c", false, "dumps configuration")

	flag.Parse()

	if *versionFlag {
		log.Println("Stor Data", version.GetVersion())
		os.Exit(0)
	}

	config.LoadConfig()

	if *configFlag {
		log.Printf("HTTP:\t%+v\n", config.HttpServer)
		log.Printf("Cassandra:\t%+v\n", config.Cassandra)

		os.Exit(0)
	}

	_, err := logger.NewRotatingFileWriter(config.App.Logfile, 500)

	if err == nil {
		log.SetOutput(logger.Logger)
	}
}

func main() {
	db, err := model.NewDb(config.Cassandra)

	if err != nil {
		log.Panic(err)
	}

	if db == nil {
		log.Panic("cassa", err)
	}

	env := &Env{db}

	log.Println("Stor Data", version.GetVersion())

	router := mux.NewRouter()

	router.HandleFunc("/data/{nodeId:[0-9]+}", env.handleDataGet).Methods("GET").Queries("start", "", "end", "")
	router.HandleFunc("/data/{nodeId:[0-9]+}/last", env.handleDataGetLast).Methods("GET")
	router.HandleFunc("/data/{nodeId:[0-9]+}", env.handleDataPost).Methods("POST")
	router.HandleFunc("/data/validate", env.handleDataValidate).Methods("POST")

	router.HandleFunc("/health", env.handleHealth)

	n := negroni.New()

	// Convert panics to 500 responses
	n.Use(negroni.NewRecovery())

	// Pretty print REST requests
	//n.Use(negroni.NewLogger())

	n.UseHandler(router)

	addr := fmt.Sprintf("%s:%d", config.HttpServer.Uri, config.HttpServer.Port)

	serv := &http.Server{
		Addr:           addr,
		Handler:        n,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Binding HTTP on", addr)

	log.Fatal("http serv:", serv.ListenAndServe())
}

func (env *Env) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (env *Env) handleDataGet(w http.ResponseWriter, r *http.Request) {
	request := model.NewDataRequest(0, time.Now(), time.Now())

	err := request.FromQuery(r.URL.Query())

	if err != nil {
		log.Println("stor-data-get query:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := request.Validate(); err != nil {
		log.Println("stor-data-get validation:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%+v", request)

	result, err := env.db.GetData(request)

	if err != nil {
		log.Println("stor-data-get exec:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(result)
}

func (env *Env) handleDataGetLast(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	nodeId, err := strconv.Atoi(vars["nodeId"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data, err := env.db.GetLast(nodeId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(data)
}

func (env *Env) handleDataPost(w http.ResponseWriter, r *http.Request) {
	data := model.Data{}

	if err := data.FromJson(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := data.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := env.db.SaveData(&data)

	if err != nil {
		log.Println("POST data:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (env *Env) handleDataValidate(w http.ResponseWriter, r *http.Request) {
	data := model.Data{}

	if err := data.FromJson(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("%+v\n", data)

	if err := data.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}