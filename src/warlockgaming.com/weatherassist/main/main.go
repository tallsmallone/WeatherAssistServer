// main.go
// Author: Nick Flint
// Date: 08/09/2017
// Description:  Main service for WeatherAssist

package main

import (
	"net/http"
	"bytes"
	"warlockgaming.com/weatherassist/config"
	"fmt"
	"io/ioutil"
	)
	
// global url
var api_url string
const max_calls_per_day int = 500
const max_calls_per_ten_minutes int = 10

func setLocation(location string) string {
	var buffer bytes.Buffer
	
	// create query string
	buffer.WriteString(api_url)
	buffer.WriteString(location)
	buffer.WriteString(".json")
	
	return buffer.String()
}

func displayResponse(w http.ResponseWriter, req *http.Request) {
	var page string
	
	url := setLocation(req.URL.Path[1:])
	resp, err := http.Get(url)
	
	if err != nil {
		page = "Error loading page"
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		
		if err != nil {
			page = "Error loading page"
		} else {
			page = string(body[:])
		}
	}
	//page = url
	fmt.Fprintf(w, page)
}

func initialize() {
	var buffer bytes.Buffer
	
	// create query string
	buffer.WriteString("http://api.wunderground.com/api/")
	buffer.WriteString(config.API_KEY)
	buffer.WriteString("/hourly/q/")
	api_url = buffer.String()
	
	// create handlers
	http.HandleFunc("/", displayResponse)
	
	// statically import all public assests
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	
	// setup listening port
	http.ListenAndServe(":8080", nil)
}

func main() {
	initialize()
}