package Handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

func SearchHandler(w http.ResponseWriter, r *http.Request, artists []FullData) {
	templates, err = template.ParseFiles(
		"templates/search.html")
		if err != nil {
			InternalServerErrorHandler(w,r)
			return
		}
		
	if r.Method != http.MethodGet {
		MethodNotAllowedHandler(w,r)
		return
	}
	text := r.FormValue("text")
	if len(text) == 0 {
		BadRequestHandler(w,r)
		return
	}
	result, err := FindData(strings.TrimSpace(text), artists)
	if err != nil {
		BadRequestHandler(w,r)
		return
	}
	HomePageHandler(w,r, result)
}

func FindData(text string, allData []FullData) ([]FullData, error) {
	var result []FullData

	text = strings.ToLower(text)
	for _, v := range allData {
		if strings.Contains(strings.ToLower(v.Name), text) {
			result = append(result, v)
			continue
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
		// myError := errors.New(http.StatusBadRequest)
		myError:=err
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