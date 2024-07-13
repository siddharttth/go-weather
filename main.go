package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

const openWeatherApiKey = "f783234b54da68d3d900b689f2aa2808"

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func query(city string) (weatherData, error) {
	if openWeatherApiKey == "" {
		return weatherData{}, fmt.Errorf("API key is missing")
	}

	url := "http://api.openweathermap.org/data/2.5/weather?APPID=" + openWeatherApiKey + "&q=" + city
	fmt.Println("Request URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		return weatherData{}, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return weatherData{}, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var d weatherData
	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return weatherData{}, fmt.Errorf("failed to decode response: %v", err)
	}

	fmt.Printf("Response data: %+v\n", d)
	return d, nil
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City name is required", http.StatusBadRequest)
		return
	}

	data, err := query(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/weather.html?city=%s&temp=%f", data.Name, data.Main.Kelvin), http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/weather/", weatherHandler)
	http.Handle("/weather.html", http.FileServer(http.Dir(".")))
	http.Handle("/styles.css", http.FileServer(http.Dir(".")))

	fmt.Println("Server started at :3000")
	http.ListenAndServe(":3000", nil)
}
