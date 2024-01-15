package Handlers

import (
	"html/template"
	"net/http"
)

func HomePageHandler(w http.ResponseWriter, r *http.Request, artists []FullData) {
	templates, err = template.ParseFiles(
		"templates/index.html")
		if err != nil {
			InternalServerErrorHandler(w,r)
			return
		}
if r.URL.Path == "/" {

err := templates.ExecuteTemplate(w, "index.html", artists)
if err != nil {
	InternalServerErrorHandler(w,r)
}
}else{
	NotFoundHandler(w,r)
	return
}
	// fmt.Printf("Data: %+v\n", ArtistsFull)
}
