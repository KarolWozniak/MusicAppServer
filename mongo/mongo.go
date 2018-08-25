package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Song struct {
	Title          string
	DownloadNumber int
	DownloadUrl    string
}

func SaveInDatabase(s *mgo.Session, songName string, url string) {
	session := s.Copy()
	defer session.Close()
	c := session.DB("test").C("songs")
	result := Song{}
	err := c.Find(bson.M{"title": songName}).One(&result)
	if result.Title != "" {
		result.DownloadNumber = result.DownloadNumber + 1
		err = c.Update(bson.M{"title": songName}, &result)
	} else {
		err = c.Insert(&Song{songName, 1, url})
	}
	if err != nil {
		log.Fatal(err)
		return
	}
}

func GetRankingFromDatabase(s *mgo.Session) []Song {
	session := s.Copy()
	defer session.Close()
	c := session.DB("test").C("songs")
	var result []Song
	ranking := c.Find(bson.M{}).Sort("-downloadnumber").Limit(3).Iter()
	err := ranking.All(&result)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return result
}
