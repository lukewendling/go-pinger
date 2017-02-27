package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"

	"github.com/gin-gonic/gin"
)

type Stats struct {
	RespTime   int     `form:"resp_time" json:"resp_time"`
	EventCount float64 `form:"event_count" json:"event_count"`
	Created    int     `form:"created" json:"created"`
}

var db *mgo.Session

func main() {
	var err error

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "mongo:27017"
	}

	db, err = mgo.Dial(dbHost)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Optional. Switch the session to a monotonic behavior.
	// dbSession.SetMode(mgo.Monotonic, true)

	router := gin.Default()

	router.Static("/app", "../app")
	// router.StaticFS("/more_static", http.Dir("my_file_system"))
	// router.StaticFile("/favicon.ico", "./resources/favicon.ico")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/app")
	})

	router.GET("/test/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.POST("/api/watchman", func(c *gin.Context) {
		var json Stats
		if err := c.BindJSON(&json); err == nil {
			fmt.Printf("%+v\n", json)

			saveToDb(json)
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		} else {
			fmt.Println(err)

		}
	})

	router.GET("/api/watchman", func(c *gin.Context) {

		results := []Stats{}

		coll := db.DB("app_stats").C("events")
		err := coll.Find(nil).Limit(3600).Sort("-created").All(&results)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, &results)

	})

	router.Run()
}

func saveToDb(stats Stats) {

	coll := db.DB("app_stats").C("events")
	err := coll.Insert(stats)
	if err != nil {
		log.Fatal(err)
	}

}
