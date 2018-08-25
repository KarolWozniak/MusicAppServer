package api

import (
	"awesomeProject/mongo"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"net/http"
	"os/exec"
)

type Video struct {
	Title       string `json:"title"`
	DownloadURL string `json:"downloadURL"`
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
			go mongo.SaveInDatabase(s, response.Title, response.DownloadURL)
			json.NewEncoder(w).Encode(response)
		}
	}
}

func GetRanking(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if result := mongo.GetRankingFromDatabase(s); result != nil {
			json.NewEncoder(w).Encode(result)
		}
	}
}
