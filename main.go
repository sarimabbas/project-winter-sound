package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"regexp"
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

// tests a string to make sure it is a valid email address using regex
func isValidEmail(toTest string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(toTest) > 254 || !rxEmail.MatchString(toTest) {
		return false
	}
	return true
}

func main() {
	// gets url of db from env variables
	dbUrl := os.Getenv("DBURL")

	// opens a new connection to db using gorm package
	db, err := gorm.Open("postgres", dbUrl)

	if err != nil {
		panic("failed to connect database")
	}

	defer db.Close()

	// db set up if models aren't already connected
	db.AutoMigrate(&Event{}, &RSVP{})

	r := gin.Default()

	// loads html templates
	r.LoadHTMLFiles("html/index.html",
		"html/about.html",
		"html/event.html",
		"html/new.html",
		"html/components.html")

	if err != nil {
		panic("failed to load html files")
	}

	// loads static assets
	r.Static("/assets", "./assets")

	// handler for home route
	r.GET("/", func(c *gin.Context) {
		events := make([]Event, 0)
		db.Find(&events)
		c.HTML(200, "index.html", gin.H{
			"title":  "Main website",
			"today":  time.Now(),
			"events": events,
		})
	})

	// handler for about page
	r.GET("/about", func(c *gin.Context) {
		c.HTML(200, "about.html", gin.H{
			"title": "About page",
		})
	})

	// handler for events page
	r.GET("/api/events", func(c *gin.Context) {
		events := make([]Event, 0)
		db.Find(&events)
		c.JSON(200, gin.H{
			"events": events,
		})
	})

	// handler for events api
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

	// handler for events detail page
	r.GET("/events/:id", func(c *gin.Context) {
		// get id
		id := c.Param("id")
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		supportType := r1.Intn(100)

		if supportType < 50 {
			supportType = 0
		} else {
			supportType = 1
		}
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

			formattedDate := event.Date.Format("2006-01-02T15:04:05")
			formattedDate = formattedDate[0 : len(formattedDate)-3]

			c.HTML(200, "event.html", gin.H{
				"event":            event,
				"rsvps":            rsvps,
				"formattedDate":    formattedDate,
				"showConfirm":      showConfirm,
				"confirmationCode": confirmationCode,
				"rsvpError":        showRSVPError,
				"supportType":      supportType,
			})
		}
	})

	// handler for new events post route
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

		log.Println(date)

		datetime := strings.Split(date, "T")

		log.Println(datetime)

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

		// e.g. 2019-12-17
		dateStr := strings.Split(datetime[0], "-")

		log.Println(dateStr)

		// e.g. [2019, 12, 17]
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

		// e.g. 2019
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

		// e.g. 2019
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

		// e.g. 12
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

		// e.g. 12
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

		// e.g. 17, 01, 05 etc
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

		// e.g. 17, 01, 05 etc
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

		log.Println(timeStr)

		// if len(timeStr) != 2 {
		// 	c.HTML(200, "new.html", gin.H{
		// 		"errorDatetime": "Invalid Date",
		// 		"eventTitle":    title,
		// 		"eventLocation": location,
		// 		"eventImage":    image,
		// 		"eventDate":     date,
		// 	})
		// 	return
		// }

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

		c.Redirect(302, "/events/"+strconv.Itoa(event.ID))
		c.Abort()

	})

	r.POST("/rsvp_events/:id", func(c *gin.Context) {
		id := c.Param("id")
		idNum, err := strconv.Atoi(id)
		email := c.PostForm("email")
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		supportType := r1.Intn(100)
		if supportType < 50 {
			// donate
			supportType = 0
		} else {
			// support
			supportType = 1
		}
		if !isValidEmail(email) {
			c.Redirect(301, "/events/"+id+"?rsvp_error="+"true")
			return
		}

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

			//c.Redirect(301, "/events/"+id+"?confirmation="+hashString[:7])
			//c.Abort()
			event := Event{}
			rsvps := make([]RSVP, 0)
			db.Where("id = ?", idNum).First(&event)
			db.Where("event_id = ?", idNum).Find(&rsvps)
			c.HTML(200, "event.html", gin.H{
				"event":            event,
				"rsvps":            rsvps,
				"showConfirm":      true,
				"confirmationCode": hashString[:7],
				"rsvpError":        false,
				"supportType":      supportType,
			})
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
