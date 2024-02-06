package Handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "log"
	"net/http"
	"time"
)

func FetchData(url string, data interface{}) error {
	client := &http.Client{
		Timeout: 10 * time.Second, // Set a timeout for the request
	}

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("Request failed with status code: %d", response.StatusCode)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
	}

	// fmt.Printf("Data: %+v\n", data)
	return nil
}
