package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type CurrentWeather struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		Condition struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
		} `json:"condition"`
	} `json:"current"`
}

func getWeatherData(city string, apiKey string) (*CurrentWeather, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no",
		apiKey, city)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var data CurrentWeather
	if err := json.NewDecoder(response.Body).Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	apiKey := "1815bf5c83254402bdb141029212309"
	// City from the Form
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}
	city := r.PostFormValue("searchCity")

	// Getting data from json
	data, err := getWeatherData(city, apiKey)
	if err != nil {
		panic(err)
	}

	// connecting with html template
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err)
	}
	tmpl.Execute(w, data)
}

func main() {
	fs := http.FileServer(http.Dir("./res/"))
	http.Handle("/res/", http.StripPrefix("/res", fs))

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
