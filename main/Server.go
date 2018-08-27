package main

import (
	"awesomeProject/api"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	router := mux.NewRouter()
	router.HandleFunc("/api/converter", api.GetVideo(session)).Methods("GET")
	router.HandleFunc("/api/ranking", api.GetRanking(session)).Methods("GET")
	router.HandleFunc("/api/trends", api.GetTrends).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
