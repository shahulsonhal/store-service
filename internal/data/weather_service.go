package data

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

func getAccuweather(country string, lat, lng float64) (*Accuweather, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/accuweather/last_minute", os.Getenv("WEATHER_URL")), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("country", country)
	q.Add("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Add("lng", strconv.FormatFloat(lng, 'f', -1, 64))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{
		// TODO: set timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var accuweather Accuweather
	err = json.NewDecoder(resp.Body).Decode(&accuweather)
	if err != nil {
		return nil, err
	}

	return &accuweather, nil
}

func getAerisweather(country string, lat, lng float64) (*Aerisweather, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/aerisweather/api/v1/current", os.Getenv("WEATHER_URL")), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("country", country)
	q.Add("lat", strconv.FormatFloat(lat, 'f', -1, 64))
	q.Add("lng", strconv.FormatFloat(lng, 'f', -1, 64))
	req.URL.RawQuery = q.Encode()

	client := &http.Client{
		// TODO: set timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var aerisweather Aerisweather
	err = json.NewDecoder(resp.Body).Decode(&aerisweather)
	if err != nil {
		return nil, err
	}

	return &aerisweather, nil
}
