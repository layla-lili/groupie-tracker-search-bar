package main

import (
	"groupie/Handlers"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	urlArtists   = "https://groupietrackers.herokuapp.com/api/artists"
	urlLocations = "https://groupietrackers.herokuapp.com/api/locations"
	urlDates     = "https://groupietrackers.herokuapp.com/api/dates"
	urlRelation  = "https://groupietrackers.herokuapp.com/api/relation"
)

var fetched = true
var artists = Handlers.Artists{}
var locations = Handlers.Locations{}
var dates = Handlers.Dates{}
var Relation = Handlers.RelationsData{}

func main() {
	var wg sync.WaitGroup
	// Create a channel to receive a signal when all goroutines are done
	done := make(chan bool)
	// Fetch data concurrently using goroutines
	wg.Add(4)
	go fetchDataConcurrently(urlArtists, &artists, &wg)
	go fetchDataConcurrently(urlLocations, &locations, &wg)
	go fetchDataConcurrently(urlDates, &dates, &wg)
	go fetchDataConcurrently(urlRelation, &Relation, &wg)
	// Start a timer to wait for a certain duration before signaling the completion
	go func() {
		// Wait for 5 seconds
		time.Sleep(5 * time.Second)
		done <- true
	}()
	// Wait for all goroutines to complete or the timer to expire
	go func() {
		wg.Wait()
		done <- true
	}()

	// Wait for the completion signal
	<-done
	staticDir := http.Dir("static")
	fs := http.FileServer(staticDir)
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	// Register the handlers
	http.HandleFunc("/", handleRequest)
	http.HandleFunc("/details", handleRequest)
	http.HandleFunc("/search", handleRequest)
	http.HandleFunc("/404", Handlers.NotFoundHandler)
	http.HandleFunc("/400", Handlers.BadRequestHandler)
	http.HandleFunc("/405", Handlers.MethodNotAllowedHandler)
	http.HandleFunc("/500", Handlers.InternalServerErrorHandler)
	// Start the server
	log.Println("Server listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server error:", err)
	}
}

func fetchDataConcurrently(url string, data interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	startTime := time.Now()
	if Handlers.FetchData(url, data) != nil {
		fetched = false
	}
	duration := time.Since(startTime)
	log.Println("Goroutine completed in", duration)
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	if fetched {
		var ArtistsFull []Handlers.FullData // Declare ArtistsFull here
		for i := range artists {
			tmpl := Handlers.FullData{
				ID:             artists[i].ID,
				Image:          artists[i].Image,
				Name:           artists[i].Name,
				Members:        make(map[string]string),
				CreationDate:   artists[i].CreationDate,
				FirstAlbum:     artists[i].FirstAlbum,
				Locations:      locations.Index[i].Locations,
				Dates:          dates.Index[i].Dates,
				DatesLocations: Relation.Index[i].DatesLocations,
			}
			for _, member := range artists[i].Members {
				// Set the member name as both the key and value in the map
				tmpl.Members[member] = member
			}
			if tmpl.Image == "https://groupietrackers.herokuapp.com/api/images/mamonasassassinas.jpeg" {
				tmpl.Image = "static/Images/ops.jpg"
			}
			ArtistsFull = append(ArtistsFull, tmpl)
		}
		path := r.URL.Path
		switch path {
		case "/":
			Handlers.HomePageHandler(w, r, ArtistsFull)
		case "/search":
			Handlers.SearchHandler(w, r, ArtistsFull)
		case "/details":
			Handlers.DetailspageHandler(w, r, ArtistsFull)
		default:
			Handlers.NotFoundHandler(w, r)
		}
	} else {
		Handlers.InternalServerErrorHandler(w, r)
	}
}