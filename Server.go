package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"os/exec"
	"encoding/json"
)

type Video struct{
	Title string `json:"title"`
	DownloadURL string `json:"downloadURL"`
}

func runCommand(url string) Video{
	title, _:=exec.Command("youtube-dl","--get-title", url).CombinedOutput()
	downloadURL, _:= exec.Command("youtube-dl", "--get-url", "-f", "bestaudio", url).CombinedOutput()
	return Video{string(title),string(downloadURL)}
}

func GetVideo(w http.ResponseWriter, r *http.Request){
	param := r.URL.Query().Get("url")
	if param != "" {
		response:= runCommand(param)
		json.NewEncoder(w).Encode(response)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/converter", GetVideo).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
