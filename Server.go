package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"os/exec"
)

type Video struct {
	Title       string `json:"title"`
	DownloadURL string `json:"downloadURL"`
}

type Song struct {
	Title          string
	DownloadNumber int
}

func saveInDatabase(s *mgo.Session, songName string) {
	session := s.Copy()
	defer session.Close()
	c := session.DB("test").C("songs")
	result := Song{}
	err := c.Find(bson.M{"title": songName}).One(&result)
	if result.Title != "" {
		result.DownloadNumber = result.DownloadNumber + 1
		err = c.Update(bson.M{"title": songName}, &result)
	} else {
		err = c.Insert(&Song{songName, 1})
	}
	if err != nil {
		panic(err)
		return
	}
}

func getFromDatabase(s *mgo.Session) []Song {
	session := s.Copy()
	defer session.Close()
	c := session.DB("test").C("songs")
	var result []Song
	ranking := c.Find(bson.M{}).Sort("-downloadnumber").Limit(3).Iter()
	err := ranking.All(&result)
	if err != nil {
		panic(err)
		return nil
	}
	return result
}

func runCommand(url string) Video {
	title, _ := exec.Command("youtube-dl", "--get-title", url).CombinedOutput()
	downloadURL, _ := exec.Command("youtube-dl", "--get-url", "-f", "bestaudio", url).CombinedOutput()
	return Video{string(title), string(downloadURL)}
}

func GetVideo(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		param := r.URL.Query().Get("url")
		if param != "" {
			response := runCommand(param)
			saveInDatabase(s, response.Title)
			json.NewEncoder(w).Encode(response)
		}
	}
}

func GetRanking(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if result := getFromDatabase(s); result != nil {
			json.NewEncoder(w).Encode(result)
		}
	}
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()
	router := mux.NewRouter()
	router.HandleFunc("/api/converter", GetVideo(session)).Methods("GET")
	router.HandleFunc("/api/ranking", GetRanking(session)).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
