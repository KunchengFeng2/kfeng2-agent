package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"io/ioutil"
)

// The actual response contains tons more information, I'm only using what I need.
// Golang assign values by same variable names, no need to mimic the complete structure.
type ServerStatus struct {
	IP string
	Version string
	Online bool
	Hostname string
	Players struct {
		Online int
		Max int
		List []string
	}
}

func main() {
	// Making the get request, no API key required.
	var apiBaseJava = "https://api.mcsrvstat.us/2/"
	// var apiBaseBedrock = "https://api.mcsrvstat.us/bedrock/2/"
	var serverName = "play.ml-mc.com"
	var response, err = http.Get(apiBaseJava + serverName)
	if err != nil {
		panic(err)
	}

	// Read the response body
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		fmt.Println(error)
	}
	defer response.Body.Close()

	// Unmarshal the json object
	var server ServerStatus
	json.Unmarshal([]byte(body), &server)

	// Print out stuff
	fmt.Println("Response status: ", response.Status)
	fmt.Println("Hostname: ", server.Hostname)
	fmt.Println("IP: ", server.IP)
	fmt.Println("Version: ", server.Version)
	fmt.Println("Online: ", server.Online)
	fmt.Println("Max players: ", server.Players.Max)
	fmt.Println("Online players: ", server.Players.Online)
	// fmt.Println("Players list: ", server.Players.List)

	// Organize information for loggly
	var logglyURL = "http://logs-01.loggly.com/inputs/5e085983-7ed1-4fc1-bf95-5f6278278035/tag/http/"
	var data = url.Values{
		"Server name": {serverName},
		"Status": {response.Status},
	}

	// Send information to loggly
	response, err = http.PostForm(logglyURL, data)
	if err != nil {
		panic(err)
	}
}