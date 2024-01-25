package Handlers

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request, artists []FullData) {
	if r.URL.Path == "/search" {
		templates, err := template.ParseFiles("templates/search.html")
		if err != nil {
			InternalServerErrorHandler(w, r)
			return
		}
		if r.Method != http.MethodGet {
			MethodNotAllowedHandler(w, r)
			return
		}
		text := r.URL.Query().Get("text")
		result, findErr := FindData(strings.TrimSpace(text), artists)
	    
		if findErr != nil {
			BadRequestHandler(w, r)
			return
		}
		err = templates.ExecuteTemplate(w, "search.html", result )
		
		if err != nil {
			InternalServerErrorHandler(w, r)
		}
	} else {
		NotFoundHandler(w, r)
		return
	}
}

func FindData(text string, allData []FullData) ([]FullData, error) {
	var result []FullData
	text = strings.ToLower(text)
	for _, v := range allData {
		if strings.Contains(strings.ToLower(v.Name), text) {
			if Check(v.ID, result) {
			result = append(result, v)
			continue
			}
		}
		if strings.Contains(v.FirstAlbum, text) {
			if Check(v.ID, result) {
				result = append(result, v)
				continue
			}
		}
		if strings.Contains(strconv.Itoa(v.CreationDate), text) {
			if Check(v.ID, result) {
				result = append(result, v)
				continue
			}
		}
		for _, member := range v.Members {
			if strings.Contains(strings.ToLower(member), text) {
				if Check(v.ID, result) {
					result = append(result, v)
				}
			}
		}
		for key := range v.DatesLocations {
			if strings.Contains(strings.ToLower(key), text) {
				if Check(v.ID, result) {
					result = append(result, v)
				}
			}
		}
	}
	if result == nil {
		myError := errors.New("BadRequest")
		return nil, myError
	}
	return result, nil
}

func Check(id int, Data []FullData) bool {
	for _, m := range Data {
		if m.ID == id {
			return false
		}
	}
	return true
}