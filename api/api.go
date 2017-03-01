package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	mgo "gopkg.in/mgo.v2"
	_ "gopkg.in/mgo.v2/bson"

	"path/filepath"
	"time"

	"gopkg.in/gin-gonic/gin.v1"
)

type Stats struct {
	QueryType string  `form:"query_type" json:"query_type"`
	RespTime  int     `form:"resp_time" json:"resp_time"`
	Count     float64 `form:"count" json:"count"`
	Created   int     `form:"created" json:"created"`
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

	fp, err := filepath.Abs("../app/index.tmpl")
	if err != nil {
		panic(err)
	}

	router.LoadHTMLFiles(fp)

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/dash")
	})

	router.GET("/dash", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"cache_bust": time.Now().Nanosecond() / 1000000,
		})
	})

	router.Static("/app", "../app")
	router.StaticFile("/favicon.ico", "../app/favicon.ico")

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

		coll := db.DB("app_stats").C("app")
		err := coll.Find(nil).Limit(3600).Sort("-created").All(&results)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, &results)
	})

	// allow get request for ease of use
	router.GET("/api/watchman/drop", func(c *gin.Context) {
		coll := db.DB("app_stats").C("app")
		err := coll.DropCollection()
		if err != nil {
			fmt.Println(err)
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	err = router.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func saveToDb(stats Stats) {
	coll := db.DB("app_stats").C("app")
	err := coll.Insert(stats)
	if err != nil {
		log.Fatal(err)
	}
}
