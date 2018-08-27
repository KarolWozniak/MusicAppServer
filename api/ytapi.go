package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Element struct {
	Kind          string `json:"kind"`
	Etag          string `json:"etag"`
	NextPageToken string `json:"nextPageToken"`
	RegionCode    string `json:"regionCode"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
	Items []struct {
		Kind string `json:"kind"`
		Etag string `json:"etag"`
		ID   struct {
			Kind    string `json:"kind"`
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			PublishedAt time.Time `json:"publishedAt"`
			ChannelID   string    `json:"channelId"`
			Title       string    `json:"title"`
			Description string    `json:"description"`
			Thumbnails  struct {
				Default struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"default"`
				Medium struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"medium"`
				High struct {
					URL    string `json:"url"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"high"`
			} `json:"thumbnails"`
			ChannelTitle         string `json:"channelTitle"`
			LiveBroadcastContent string `json:"liveBroadcastContent"`
		} `json:"snippet"`
	} `json:"items"`
}

const apiUrl = "https://www.googleapis.com/youtube/v3/search?part=snippet&type=video&videoCategoryId=10&regionCode=US&key=API_KEY"
const ytUrl = "https://www.youtube.com/watch?v="

func GetYtRanking() []Video {
	response, err := http.Get(apiUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	var decoded Element
	err = json.NewDecoder(response.Body).Decode(&decoded)
	if err != nil {
		log.Fatal(err)
	}
	var videos []Video
	c := make(chan Video)
	for _, a := range decoded.Items {
		fmt.Println(createUrl(a.ID.VideoID), a.Snippet.Title)
		go getVideo(createUrl(a.ID.VideoID), c)
	}
	for range decoded.Items {
		videos = append(videos, <-c)
	}
	return videos
}

func createUrl(videoId string) string {
	return ytUrl + videoId
}

func getVideo(url string, c chan Video) {
	result := runCommand(url)
	c <- result
}
