package main

import (
	"net/http"
	"time"
)

func indexController(w http.ResponseWriter, r *http.Request) {

	type indexContextData struct {
		Events []Event
		Today  time.Time
	}

	allEvents, err := getAllEvents()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	contextData := indexContextData{
		Events: allEvents,
		Today:  time.Now(),
	}

	tmpl["index"].Execute(w, contextData)
}
