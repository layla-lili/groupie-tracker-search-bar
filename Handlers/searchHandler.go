package Handlers

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"sort"

)
// Define a custom sorter type
type ByName []FullData

// Implement the sort.Interface methods for ByName type
func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return strings.Compare(a[i].Name, a[j].Name) < 0 }



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
		err = templates.ExecuteTemplate(w, "search.html", result)

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
		for _, date := range v.Dates {
			if strings.Contains(strings.ToLower(date), text) {
				if Check(v.ID, result) {
					result = append(result, v)
				}
			}
		}
		for _, location := range v.Locations {
			if strings.Contains(strings.ToLower(location), text) {
				if Check(v.ID, result) {
					result = append(result, v)
				}
			}
		}
		for key := range v.DatesLocations {
			concertDates := strings.Join(v.DatesLocations[key], " ")
			for range key {
				if strings.Contains(strings.ToLower(concertDates), string(text)) {
					if Check(v.ID, result) {
						result = append(result, v)
					}
				}
			}

		}

	}
	// Sort the result slice by name
	sort.Sort(ByName(result))

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
