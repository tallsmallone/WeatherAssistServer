// main.go
// Author: Nick Flint
// Date: 08/09/2017
// Description:  Main service for WeatherAssist

package main

import (
	"net/http"
	"bytes"
	"warlockgaming.com/weatherassist/config"
	"io/ioutil"
	"time"
	"strconv"
	"strings"
	"google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
	)
	
// global url
var api_url string

// max api calls
const max_calls_per_day int = 500
const max_calls_per_minute int = 10
var calls_per_minute int = 0
var calls_per_day int = 0

// time 
var one_minute_start time.Time
var one_day_start time.Time

func getLocationText(url string) []string {
	s := strings.Split(url, "/")
	
	return s
}

func checkFirstTime() {
	
	if one_minute_start.IsZero() {
		one_minute_start = time.Now()
	}
	
	if one_day_start.IsZero() {
		one_day_start = time.Now()
	}
}

func checkTime() bool {
	minute_flag := false
	day_flag := false
	
	checkFirstTime()
	
	// check for number of attempts per 1 minute
	difference_minute := time.Since(one_minute_start)
	
	if difference_minute.Minutes() > 1.0 {
		calls_per_minute = 1
		one_minute_start = time.Now()
		minute_flag = true
	} else if calls_per_minute < max_calls_per_minute {
		calls_per_minute++
		minute_flag = true
	} else {
		minute_flag = false
	}
	
	// check for number of attempts per day
	difference_day := time.Since(one_day_start)
	if difference_day.Hours() > 24.0 {
		calls_per_day = 1
		one_day_start = time.Now()
		day_flag = true
	} else if calls_per_day < max_calls_per_day {
		calls_per_day++
		day_flag = true
	} else {
		day_flag = false
	}
	
	return minute_flag && day_flag
}

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
	
	// check for blank to prevent bad calls
	location := getLocationText(req.URL.Path[1:])[1]
	if location == string("") {
		w.Write([]byte("No location specified"))
		return
	}
	
	if !checkTime() {
		w.Write([]byte("Too many calls"))
		return
	} 

	url := setLocation(string(location[:]))
	ctx := appengine.NewContext(req)
	client := urlfetch.Client(ctx)
	resp, err := client.Get(url)
	
	//resp, err := http.Get(url)
	
	if err != nil {
		page = "Error getting page"
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
	w.Write([]byte(page))
}

func checkUpResponse(w http.ResponseWriter, req *http.Request) {
	minute := "Calls per minute: " + strconv.Itoa(calls_per_minute) + "\n"
	day := "Calls per day: " + strconv.Itoa(calls_per_day)
	w.Write([]byte(minute))
	w.Write([]byte(day))
}

func initialize() {
	var buffer bytes.Buffer
	
	// create query string
	buffer.WriteString("http://api.wunderground.com/api/")
	buffer.WriteString(config.API_KEY)
	buffer.WriteString("/hourly/q/")
	api_url = buffer.String()
	
	// create handlers
	http.HandleFunc("/query/", displayResponse)
	http.HandleFunc("/checkup/", checkUpResponse)
	
	// statically import all public assests
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))
	
	// setup listening port
	http.ListenAndServe(":8080", nil)
}

func init() {
	initialize()
}