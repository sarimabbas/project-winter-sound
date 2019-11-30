package main

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// isValidUrl tests a string to determine if it is a well-structured url or not.
func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	} else {
		return true
	}
}

func randomString() string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	length := 8
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := b.String()
	return str
}

func main() {
	dbUrl := os.Getenv("DBURL")

	db, err := gorm.Open("postgres", dbUrl)

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	db.AutoMigrate(&Event{}, &RSVP{})

	r := gin.Default()

	r.LoadHTMLFiles("html/index.html",
		"html/about.html",
		"html/event.html",
		"html/new.html",
		"html/components.html")

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

		confirmationCode := c.Query("confirmation")
		showConfirm := false
		showRSVPError := false

		if confirmationCode != "" {
			showConfirm = true
		}

		rsvpError := c.Query("rsvp_error")

		if rsvpError != "" {
			showRSVPError = true
		}
		// new event page
		if id == "new" {
			c.HTML(200, "new.html", gin.H{
				"title": "Create event page",
				"error": false,
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
			rsvps := make([]RSVP, 0)
			db.Where("id = ?", idNum).First(&event)
			db.Where("event_id = ?", idNum).Find(&rsvps)
			c.HTML(200, "event.html", gin.H{
				"event":            event,
				"rsvps":            rsvps,
				"showConfirm":      showConfirm,
				"confirmationCode": confirmationCode,
				"rsvpError":        showRSVPError,
			})
		}
	})

	r.POST("/events/new", func(c *gin.Context) {
		// get data from form
		title := c.PostForm("title")
		location := c.PostForm("location")
		image := c.PostForm("image")
		date := c.PostForm("date")

		// next steps:

		// 1.
		// validate data and return error code if invalid
		// this error is displayed on the form
		// validate image URL (has to be URL + end in image extension)
		// validate datetime string

		if len(title) < 5 || len(title) > 50 {
			c.HTML(200, "new.html", gin.H{
				"errorTitle":    "Title must be between 5 and 50 characters.",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		if len(location) < 5 || len(location) > 50 {
			c.HTML(200, "new.html", gin.H{
				"errorLocation": "Location must be more than 5 characters and less than 50.",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		urlObj, err := url.Parse(image)

		if !isValidUrl(image) {
			c.HTML(200, "new.html", gin.H{
				"errorImage":    "Invalid URL",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		if err != nil || !urlObj.IsAbs() {
			c.HTML(200, "new.html", gin.H{
				"errorImage":    "Invalid URL",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		lastFive := image[len(image)-5 : len(image)]

		fileTypeStr := strings.Split(lastFive, ".")

		if fileTypeStr[1] != "jpg" && fileTypeStr[1] != "png" && fileTypeStr[1] != "jpeg" && fileTypeStr[1] != "gif" && fileTypeStr[1] != ".gifv" {
			c.HTML(200, "new.html", gin.H{
				"errorImage":    "Invalid Image Type (must be .png, .jpg, .jpeg, .gif or .gifv)",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		datetime := strings.Split(date, "T")

		if len(datetime) != 2 {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		dateStr := strings.Split(datetime[0], "-")

		if len(dateStr) != 3 {

			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		if len(dateStr[0]) != 4 {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		yearInt, err := strconv.Atoi(dateStr[0])
		if err != nil {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		if len(dateStr[1]) != 2 {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		monthInt, err := strconv.Atoi(dateStr[1])
		if err != nil {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		if len(dateStr[2]) != 2 {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		dayInt, err := strconv.Atoi(dateStr[2])
		if err != nil {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		timeStr := strings.Split(datetime[1], ":")

		if len(timeStr) != 2 {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		if len(timeStr[0]) != 2 {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		hourInt, err := strconv.Atoi(timeStr[0])
		if err != nil {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		if len(timeStr[1]) != 2 {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		minInt, err := strconv.Atoi(timeStr[1])
		if err != nil {
			c.HTML(200, "new.html", gin.H{
				"errorDatetime": "Invalid Date",
				"eventTitle":    title,
				"eventLocation": location,
				"eventImage":    image,
				"eventDate":     date,
			})
			return
		}

		timeObj := time.Date(yearInt, time.Month(monthInt), dayInt, hourInt, minInt, 0, 0, time.UTC)

		event := Event{
			Title:    title,
			Location: location,
			Image:    image,
			Date:     timeObj,
		}

		db.Create(&event)

		events := make([]Event, 0)
		db.Find(&events)

		c.Redirect(301, "/events/"+strconv.Itoa(event.ID))
		c.Abort()

	})

	r.POST("/rsvp_events/:id", func(c *gin.Context) {
		id := c.Param("id")
		idNum, err := strconv.Atoi(id)
		email := c.PostForm("rsvp-email")

		if !strings.Contains(email, "@yale.edu") {
			c.Redirect(301, "/events/"+id+"?rsvp_error="+"true")
			return
		}
		if err != nil {
			c.JSON(404, gin.H{
				"error": "id error",
			})
		} else {
			rsvp := RSVP{
				EventID: idNum,
				Email:   email,
			}
			db.Create(&rsvp)

			hash := sha256.Sum256([]byte(email))
			hashString := fmt.Sprintf("%x", hash)

			c.Redirect(301, "/events/"+id+"?confirmation="+hashString[:7])
			c.Abort()
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
