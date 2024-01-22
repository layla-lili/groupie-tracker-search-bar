package main

import (
	"groupie/Handlers"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)


const (
	urlArtists   = "https://groupietrackers.herokuapp.com/api/artists"
	urlLocations = "https://groupietrackers.herokuapp.com/api/locations"
	urlDates     = "https://groupietrackers.herokuapp.com/api/dates"
	urlRelation  = "https://groupietrackers.herokuapp.com/api/relation"
)

var ArtistsFull = make([]Handlers.FullData, 0, len(Handlers.Artists{}))
var fetched = true
func main() {
	artists := Handlers.Artists{}
	locations := Handlers.Locations{}
	dates := Handlers.Dates{}
	Relation := Handlers.RelationsData{}

	// Handlers.FetchData(urlArtists, &artists)
	// Handlers.FetchData(urlLocations, &locations)
	// Handlers.FetchData(urlDates, &dates)
	// Handlers.FetchData(urlRelation, &Relation)

	var wg sync.WaitGroup
	// Create a channel to receive a signal when all goroutines are done
	done := make(chan bool)

	wg.Add(4)

	// Fetch data from APIs concurrently using goroutines
	go func() {
		defer wg.Done()
		startTime := time.Now()

		if Handlers.FetchData(urlArtists, &artists) != nil{
			fetched=false
		}
		duration := time.Since(startTime)
		log.Println("Goroutine 1 completed in", duration)

	}()

	go func() {
		defer wg.Done()
		startTime := time.Now()

		if Handlers.FetchData(urlLocations, &locations)  != nil{
			fetched=false
		}

		duration := time.Since(startTime)
		log.Println("Goroutine 2 completed in", duration)

	}()

	go func() {
		defer wg.Done()
		startTime := time.Now()

		if Handlers.FetchData(urlDates, &dates)  != nil{
			fetched=false
		}
		duration := time.Since(startTime)
		log.Println("Goroutine 3 completed in", duration)

	}()

	go func() {
		defer wg.Done()
		startTime := time.Now()
		if Handlers.FetchData(urlRelation, &Relation)  != nil{
			fetched=false
		}
		duration := time.Since(startTime)
		log.Println("Goroutine 4 completed in", duration)

	}()
	
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if fetched{
		for i := range artists {

			tmpl := Handlers.FullData{
				ID:      artists[i].ID,
				Image:   artists[i].Image,
				Name:    artists[i].Name,
				Members: make(map[string]string),
	
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
		
		Handlers.HomePageHandler(w, r, ArtistsFull)
		}else{
			Handlers.InternalServerErrorHandler(w,r)
		}
	})

	http.HandleFunc("/details", func(w http.ResponseWriter, r *http.Request) {
			if fetched{
				for i := range artists {
		
					tmpl := Handlers.FullData{
						ID:      artists[i].ID,
						Image:   artists[i].Image,
						Name:    artists[i].Name,
						Members: make(map[string]string),
			
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
				
			Handlers.DetailspageHandler(w, r, ArtistsFull)
			}else{
				Handlers.InternalServerErrorHandler(w,r)
			}
	})

	// http.HandleFunc("/search", Handlers.SearchHandler)
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		if fetched{
			for i := range artists {
	
				tmpl := Handlers.FullData{
					ID:      artists[i].ID,
					Image:   artists[i].Image,
					Name:    artists[i].Name,
					Members: make(map[string]string),
		
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
			
		Handlers.SearchHandler(w, r, ArtistsFull)
		}else{
			Handlers.InternalServerErrorHandler(w,r)
		}
})

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

// Open the server URL in the default browser
openBrowser("http://localhost:8080")

	
}
// openBrowser opens the specified URL in the default browser of the user's operating system
func openBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	default:
		err = os.ErrInvalid
	}

	if err != nil {
		log.Println("Failed to open browser:", err)
	}
}
