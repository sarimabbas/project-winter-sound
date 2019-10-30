package main

import (
	"os"
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

	r.LoadHTMLFiles("html/index.html", "html/about.html")

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

	r.GET("api/events", func(c *gin.Context) {
		events := make([]Event, 0)
		db.Find(&events)
		c.JSON(200, gin.H{
			"events": events,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
