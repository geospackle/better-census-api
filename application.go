package main

import (
	"census-api/fetchdata"
	"census-api/helpers"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
    "github.com/gorilla/handlers"
	//	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	//	"time"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Read docs: https://github.com/geospackle/better-census-api")
}

func getAllDatasets() *fetchdata.Datasets {
	allDatasets := &fetchdata.Datasets{} // or &Foo{}
	helpers.GetJSON("https://api.census.gov/data/", allDatasets)
	return allDatasets
}

var allDatasets = getAllDatasets()

func findCensusDataset(w http.ResponseWriter, r *http.Request) {
	vintage := mux.Vars(r)["vintage"]
	searchTerm := mux.Vars(r)["search"]
	dataset, err := fetchdata.FindDataset(allDatasets.Dataset, vintage, searchTerm)
	if err != nil {
		http.Error(w, "error: "+err.Error(), 500)
	} else {
		out, err := json.MarshalIndent(dataset, "", "    ")
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}

func findCensusTable(w http.ResponseWriter, r *http.Request) {
	datasetID, err := strconv.Atoi(mux.Vars(r)["datasetid"])
	if err != nil {
		panic(err)
	}
	searchTerm := mux.Vars(r)["search"]
	table, err := fetchdata.FindTable(allDatasets.Dataset, datasetID, searchTerm)
	if err != nil {
		http.Error(w, "error: "+err.Error(), 500)
	} else {
		out, err := json.MarshalIndent(table, "", "    ")
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}
}

func getCensusTable(w http.ResponseWriter, r *http.Request) {
	vintage, _ := strconv.Atoi(mux.Vars(r)["vintage"])
	dataset := mux.Vars(r)["dataset"]
	group := mux.Vars(r)["group"]
	variable := mux.Vars(r)["variable"]
	geography := mux.Vars(r)["geography"]
	state := mux.Vars(r)["state"]
	county := mux.Vars(r)["county"]
	key := mux.Vars(r)["key"]
	table, statusCode := fetchdata.GetTable(key, vintage, dataset, group, variable, geography, state, county)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(table)
}

func main() {
    router := mux.NewRouter().StrictSlash(true)
    headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
    originsOk := handlers.AllowedOrigins([]string{"*"})
    methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
    router.HandleFunc("/", homeLink)
    router.HandleFunc("/finddataset", findCensusDataset).Queries("search", "{search}", "vintage", "{vintage}").Methods("GET")
    router.HandleFunc("/findtable", findCensusTable).Queries("search", "{search}", "datasetid", "{datasetid}").Methods("GET")
    router.HandleFunc("/gettable", getCensusTable).Queries("key", "{key}", "vintage", "{vintage}", "dataset", "{dataset}", "group", "{group}", "variable", "{variable}", "geography", "{geography}", "state", "{state}", "county", "{county}").Methods("GET")
    router.PathPrefix("/.well-known/pki-validation/").Handler(http.StripPrefix("/.well-known/pki-validation/",http.FileServer(http.Dir("./static/")))).Methods("GET")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS(headersOk, originsOk, methodsOk)(router)))
}
