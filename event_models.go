package main

import "time"

type Event struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Location string    `json:"location"`
	Image    string    `json:"image"`
	Date     time.Time `json:"date"`
}

type RSVP struct {
	ID      int    `json:"id"`
	EventID int    `json:"eventID"`
	Email   string `json:"email"`
}
