package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

)

func main(){
	db, err := gorm.Open("postgres", "postgres://xzewemse:KsPpTENkwPnhzYvE7kLjx1m98mXwrbjQ@salt.db.elephantsql.com:5432/xzewemse")


	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()


	db.AutoMigrate(&Event{})


	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("api/events", func(c *gin.Context){
		events := make([]Event, 0)
		db.Find(&events)
		c.JSON(200, gin.H{
			"events" : events,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}