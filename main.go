package main

import (
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	dbUrl := os.Getenv("DBURL")

	db, err := gorm.Open("postgres", dbUrl)

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&Event{})

	r := gin.Default()

	r.LoadHTMLFiles("html/index.html", "html/about.html", "html/event.html", "html/new.html")

	if err != nil {
		panic("failed to load html files")
	}

	r.Static("/assets", "./assets")

	r.GET("/", func(c *gin.Context) {
		events := make([]Event, 0)
		db.Find(&events)
		c.HTML(200, "index.html", gin.H{
			"title":  "Main website",
			"today":  time.Now(),
			"events": events,
		})
	})

	r.GET("/about", func(c *gin.Context) {
		c.HTML(200, "about.html", gin.H{
			"title": "About page",
		})
	})

	r.GET("/api/events", func(c *gin.Context) {
		events := make([]Event, 0)
		db.Find(&events)
		c.JSON(200, gin.H{
			"events": events,
		})
	})

	r.GET("/api/events/:id", func(c *gin.Context) {
		id := c.Param("id")
		idNum, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(404, gin.H{
				"error": "id error",
			})
		} else {
			event := Event{}
			db.Where("id = ?", idNum).First(&event)
			c.JSON(200, event)
		}
	})

	r.GET("/events/:id", func(c *gin.Context) {
		// get id
		id := c.Param("id")

		// new event page
		if id == "new" {
			c.HTML(200, "new.html", gin.H{
				"title": "Create event page",
			})
			return
		}

		// event detail page
		idNum, err := strconv.Atoi(id)
		if err != nil {
			c.JSON(404, gin.H{
				"error": "id error",
			})
		} else {
			event := Event{}
			db.Where("id = ?", idNum).First(&event)
			c.HTML(200, "event.html", gin.H{
				"event": event,
			})
		}
	})

	r.POST("/events/new", func(c *gin.Context) {
		// get data from form
		title := c.PostForm("event-title")
		location := c.PostForm("event-location")
		image := c.PostForm("event-image")
		date := c.PostForm("event-datetime")

		// next steps:

		// 1.
		// validate data and return error code if invalid
		// this error is displayed on the form
		// validate image URL (has to be URL + end in image extension)
		// validate datetime string

		// 2.
		// create new event in database

		// 3.
		// redirect to newly created event page

		// display (debug)
		c.JSON(200, gin.H{
			"title":    title,
			"location": location,
			"image":    image,
			"date":     date,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
